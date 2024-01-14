package main

import (
	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/cfindlayisme/rss-wmb/rss"
)

func main() {
	feedURLs := env.GetFeedUrls()
	feedChannels := env.GetFeedChannels()
	scheduledDuration := env.GetScheduledMinutes()

	rss.CheckFeeds(feedChannels, feedURLs)
	rss.ScheduleFeeds(scheduledDuration, feedChannels, feedURLs)

	// Keep the main goroutine running
	select {}
}
