package main

import (
	"time"

	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/cfindlayisme/rss-wmb/rss"
)

func main() {
	// Define the RSS feed URLs
	feedURLs := env.GetFeedUrls()

	feedChannels := env.GetFeedChannels()

	rss.CheckFeeds(feedChannels, feedURLs)
	ticker := time.NewTicker(5 * time.Minute)

	// Start a goroutine to check RSS feeds periodically
	go func() {
		for {
			select {
			case <-ticker.C:
				rss.CheckFeeds(feedChannels, feedURLs)
			}
		}
	}()

	// Keep the main goroutine running
	select {}
}
