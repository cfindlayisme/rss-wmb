package main

import (
	"time"

	"github.com/cfindlayisme/rss-wmb/db"
	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/cfindlayisme/rss-wmb/rss"
)

func main() {
	feedURLs := env.GetFeedUrls()
	feedChannels := env.GetFeedChannels()
	scheduledDuration := env.GetScheduledMinutes()

	rss.CheckFeeds(feedChannels, feedURLs)
	rss.ScheduleFeeds(scheduledDuration, feedChannels, feedURLs)

	// Once per day cleanup DB (also at start)
	go func() {
		for range time.Tick(24 * time.Hour) {
			db.CleanDB()
		}
	}()

	// Keep the main goroutine running
	select {}
}
