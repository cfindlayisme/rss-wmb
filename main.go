package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"maps"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cfindlayisme/rss-wmb/db"
	"github.com/cfindlayisme/rss-wmb/model"
	"github.com/mmcdole/gofeed"
)

func main() {
	// Define the RSS feed URLs
	feedURLs := strings.Split(os.Getenv("FEED_URLS"), ",")

	feedChannels := strings.Split(os.Getenv("FEED_CHANNELS"), ",")

	checkFeeds(feedChannels, feedURLs)
	ticker := time.NewTicker(15 * time.Minute)

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
			fmt.Printf("Error parsing RSS feed: %s\n", err)
			continue
		}

		// Create a new map to store the updated feed items
		updatedFeedItems := make(map[string]bool)

		// Print the feed items
		for _, item := range feed.Items {
			// Check if the feed item has already been printed
			if !db.GetIfLinkPrintedInDB(item.Link) {
				fmt.Printf("Title: %s\n", item.Title)
				fmt.Println("--------------------")

				// Create a new WebhookMessage
				webhookMessage := model.WebhookMessage{
					Message:  "Title: " + item.Title + " Link: " + item.Link,
					Password: os.Getenv("PASSWORD"),
				}

				webhookDirectedMessage := model.DirectedWebhookMessage{
					IncomingMessage: webhookMessage,
					Target:          feedChannels[n],
				}

				jsonData, err := json.Marshal(webhookDirectedMessage)
				if err != nil {
					log.Fatalf("Error marshalling webhookDirectedMessage: %v", err)
				}

				fmt.Printf("JSON Data: %s\n", jsonData)

				// Send a POST request to the webhook URL
				resp, err := http.Post(os.Getenv("WMB_DIRECT_MESSAGE_URL"), "application/json", bytes.NewBuffer(jsonData))
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
				fmt.Printf("Response: %s\n", body)
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
