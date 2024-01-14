package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"maps"
	"net/http"
	"time"

	"github.com/cfindlayisme/rss-wmb/db"
	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/cfindlayisme/rss-wmb/model"
	"github.com/mmcdole/gofeed"
)

func main() {
	// Define the RSS feed URLs
	feedURLs := env.GetFeedUrls()

	feedChannels := env.GetFeedChannels()

	checkFeeds(feedChannels, feedURLs)
	ticker := time.NewTicker(5 * time.Minute)

	// Start a goroutine to check RSS feeds periodically
	go func() {
		for {
			select {
			case <-ticker.C:
				checkFeeds(feedChannels, feedURLs)
			}
		}
	}()

	// Keep the main goroutine running
	select {}
}

// Function to check RSS feeds
func checkFeeds(feedChannels []string, feedURLs []string) {
	feedItemsNew := make(map[string]bool)

	// Iterate over the feed URLs
	for n, url := range feedURLs {
		// Parse the RSS feed
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err != nil {
			log.Printf("Error parsing RSS feed: %s\n", err)
			continue
		}

		// Create a new map to store the updated feed items
		updatedFeedItems := make(map[string]bool)

		// Print the feed items
		for _, item := range feed.Items {
			// Check if the feed item has already been printed
			if !db.GetIfLinkPrintedInDB(item.Link) {
				log.Printf("Title: %s\n", item.Title)
				log.Println("--------------------")

				// Create a new WebhookMessage
				webhookMessage := model.WebhookMessage{
					Message:  "Title: " + item.Title + " Link: " + item.Link,
					Password: env.GetWMBPassword(),
				}

				webhookDirectedMessage := model.DirectedWebhookMessage{
					IncomingMessage: webhookMessage,
					Target:          feedChannels[n],
				}

				jsonData, err := json.Marshal(webhookDirectedMessage)
				if err != nil {
					log.Fatalf("Error marshalling webhookDirectedMessage: %v", err)
				}

				log.Printf("JSON Data: %s\n", jsonData)

				// Send a POST request to the webhook URL
				resp, err := http.Post(env.GetWMBURL(), "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					log.Fatalf("Error sending webhook: %v", err)
				}
				// Read the response body
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("Error reading response body: %v", err)
				}

				// Print the response body
				log.Printf("Response: %s\n", body)
				defer resp.Body.Close()

				// Mark the feed item as printed
				updatedFeedItems[item.Link] = true
			}
		}
		maps.Copy(feedItemsNew, updatedFeedItems)
	}

	if len(feedItemsNew) != 0 {
		db.WriteFeedItemsToDB(feedItemsNew)
	}
}
