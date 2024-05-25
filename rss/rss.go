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

type FeedChecker interface {
	CheckFeeds(database *db.DB, feedChannels []string, feedURLs []string)
}

type Scheduler struct {
	FeedChecker FeedChecker
}

func (s *Scheduler) CheckFeeds(database *db.DB, feedChannels []string, feedURLs []string) {
	s.FeedChecker.CheckFeeds(database, feedChannels, feedURLs)
}

func (s *Scheduler) ScheduleFeeds(database *db.DB, d time.Duration, feedChannels []string, feedURLs []string) {
	if d <= 0 {
		log.Println("Error: non-positive interval for NewTicker")
		return
	}

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
	// Implement the feed checking logic here
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

func NewScheduler() *Scheduler {
	feedChecker := &DefaultFeedChecker{}
	return &Scheduler{FeedChecker: feedChecker}
}
