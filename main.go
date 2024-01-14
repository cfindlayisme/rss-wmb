package main

import (
	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/cfindlayisme/rss-wmb/rss"
)

func main() {
	// Define the RSS feed URLs
	feedURLs := env.GetFeedUrls()

	feedChannels := env.GetFeedChannels()

	rss.CheckFeeds(feedChannels, feedURLs)
	rss.ScheduleFeeds(env.GetScheduledMinutes(), feedChannels, feedURLs)

	// Keep the main goroutine running
	select {}
}
