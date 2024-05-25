package main

import (
	"log"
	"time"

	"github.com/cfindlayisme/rss-wmb/db"
	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/cfindlayisme/rss-wmb/rss"
)

func main() {
	feedURLs := env.GetFeedUrls()
	feedChannels := env.GetFeedChannels()
	scheduledDuration := env.GetScheduledMinutes()

	// Create a DB
	database, err := db.NewDB()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	scheduler := rss.NewScheduler()
	scheduler.CheckFeeds(database, feedChannels, feedURLs)
	scheduler.ScheduleFeeds(database, scheduledDuration, feedChannels, feedURLs)

	// Once per day cleanup DB (also at start)
	go func() {
		database.CleanDB()

		for range time.Tick(24 * time.Hour) {
			database.CleanDB()
		}
	}()

	// Keep the main goroutine running
	select {}
}
