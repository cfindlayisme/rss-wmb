package rss

import (
	"log"
	"maps"
	"time"

	"github.com/cfindlayisme/rss-wmb/db"
	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/cfindlayisme/rss-wmb/wmb"
	"github.com/mmcdole/gofeed"
)

func CheckFeeds(database *db.DB, feedChannels []string, feedURLs []string) {
	feedItemsNew := make(map[string]bool)

	for n, url := range feedURLs {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err != nil {
			log.Printf("Error parsing RSS feed: %s\n", err)
			continue
		}

		// Create a new map to store the updated feed items
		updatedFeedItems := make(map[string]bool)

		for _, item := range feed.Items {
			if !database.GetIfLinkPrintedInDB(item.Link) {
				wmb.SendDirectedRSSMessage(env.GetWMBURL(), item, feedChannels, n)

				// Mark the feed item as printed
				updatedFeedItems[item.Link] = true
			}
		}

		maps.Copy(feedItemsNew, updatedFeedItems)
	}

	if len(feedItemsNew) != 0 {
		database.WriteFeedItemsToDB(feedItemsNew)
	}
}

type FeedChecker interface {
	CheckFeeds(database *db.DB, feedChannels []string, feedURLs []string)
}

type Scheduler struct {
	FeedChecker FeedChecker
}

func (s *Scheduler) ScheduleFeeds(database *db.DB, d time.Duration, feedChannels []string, feedURLs []string) {
	ticker := time.NewTicker(d)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.FeedChecker.CheckFeeds(database, feedChannels, feedURLs)
			}
		}
	}()
}

type DefaultFeedChecker struct{}

func (d *DefaultFeedChecker) CheckFeeds(database *db.DB, feedChannels []string, feedURLs []string) {
	CheckFeeds(database, feedChannels, feedURLs)
}

func ScheduleFeeds(database *db.DB, d time.Duration, feedChannels []string, feedURLs []string) {
	feedChecker := &DefaultFeedChecker{}
	scheduler := &Scheduler{FeedChecker: feedChecker}
	scheduler.ScheduleFeeds(database, d, feedChannels, feedURLs)
}
