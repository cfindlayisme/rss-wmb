package env

import (
	"os"
	"strings"
)

func GetFeedUrls() []string {
	return strings.Split(os.Getenv("FEED_URLS"), ",")
}

func GetFeedChannels() []string {
	return strings.Split(os.Getenv("FEED_CHANNELS"), ",")
}
