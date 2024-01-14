package env

import (
	"os"
	"strings"
	"time"
)

func GetFeedUrls() []string {
	return strings.Split(os.Getenv("FEED_URLS"), ",")
}

func GetFeedChannels() []string {
	return strings.Split(os.Getenv("FEED_CHANNELS"), ",")
}

func GetWMBPassword() string {
	return os.Getenv("PASSWORD")
}

func GetWMBURL() string {
	return os.Getenv("WMB_DIRECT_MESSAGE_URL")
}

func GetStateFilePath() string {
	return os.Getenv("STATEFILE")
}

func GetScheduledMinutes() time.Duration {
	return 5 * time.Minute // static for now
}
