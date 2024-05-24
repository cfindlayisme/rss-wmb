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

func CheckFeeds(feedChannels []string, feedURLs []string) {
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
			if !db.GetIfLinkPrintedInDB(item.Link) {
				wmb.SendDirectedRSSMessage(env.GetWMBURL(), item, feedChannels, n)

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

func ScheduleFeeds(scheduledDuration time.Duration, feedChannels []string, feedURLs []string) {
	ticker := time.NewTicker(scheduledDuration)

	go func() {
		for {
			select {
			case <-ticker.C:
				CheckFeeds(feedChannels, feedURLs)
			}
		}
	}()
}
