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

	"github.com/cfindlayisme/rss-wmb/model"
	"github.com/mmcdole/gofeed"
)

func main() {
	// Define the RSS feed URLs
	feedURLs := strings.Split(os.Getenv("FEED_URLS"), ",")

	feedChannels := strings.Split(os.Getenv("FEED_CHANNELS"), ",")

	checkFeeds(feedChannels, feedURLs)

	// Create a ticker that ticks every ten minutes
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
	// Read the feed items from a file
	feedItems, err := readFeedItemsFromFile()
	if err != nil {
		fmt.Printf("Error reading feed items from file: %s\n", err)
		return
	}
	feedItemsNew := make(map[string]bool)

	if err != nil {
		fmt.Printf("Error reading feed items from file: %s\n", err)
	}

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
			if !feedItems[item.Link] {
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
		maps.Copy(feedItemsNew, feedItems)
		// Write the feed items to a file
		err = writeFeedItemsToFile(feedItemsNew)
		if err != nil {
			fmt.Printf("Error writing feed items to file: %s\n", err)
		}
	}
}

// Function to read the feed items from a file
func readFeedItemsFromFile() (map[string]bool, error) {
	filePath := os.Getenv("STATEFILE")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var feedItems map[string]bool
	err = json.NewDecoder(file).Decode(&feedItems)
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully read", len(feedItems), "feed items from file")

	return feedItems, nil
}

// Function to write the feed items to a file
func writeFeedItemsToFile(feedItems map[string]bool) error {
	filePath := os.Getenv("STATEFILE")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(feedItems)
	if err != nil {
		return err
	}

	fmt.Println("Successfully wrote", len(feedItems), "feed items to file")

	return nil
}
