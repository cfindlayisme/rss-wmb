Quite simply, this takes RSS feeds and checks them often, then sends the title and link to a channel on IRC through my wmb bot. 

## TODO
- Check feeds in a slightly more staggered manner
- Cleanup the sqlite db on an interval
- Flesh out this README more
- Probably a race condition in CleanupDB with the scheduled feed check in terms of the database locking - should check this out

## Enviorment Variables
- `FEED_URLS` - A comma seperated list of RSS feed urls
- `FEED_CHANNELS` - A comma seperated list of IRC channels to send the feed to, matching up exactly in position with `FEED_URLS` (ie, first feed, first item in `FEED_CHANNELS`)
- `PASSWORD` - The password to wmb
- `STATEFILE` - The sqlitedb file to store the state of the feeds in

## Why?
I no longer wanted to run any sort of RSS feed monitoring software, and tend to live in IRC. Also wanted to learn to do more things in Go, and plugin to my wmb project.