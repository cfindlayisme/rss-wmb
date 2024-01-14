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
	rss.ScheduleFeeds(5*time.Minute, feedChannels, feedURLs)

	// Keep the main goroutine running
	select {}
}
