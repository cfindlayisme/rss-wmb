package env_test

import (
	"os"
	"testing"
	"time"

	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/stretchr/testify/require"
)

func TestGetFeedUrls(t *testing.T) {
	expected := []string{"http://example.com/feed1", "http://example.com/feed2"}
	os.Setenv("FEED_URLS", "http://example.com/feed1,http://example.com/feed2")
	require.Equal(t, expected, env.GetFeedUrls())
}

func TestGetFeedChannels(t *testing.T) {
	expected := []string{"channel1", "channel2"}
	os.Setenv("FEED_CHANNELS", "channel1,channel2")
	require.Equal(t, expected, env.GetFeedChannels())
}

func TestGetWMBPassword(t *testing.T) {
	expected := "password123"
	os.Setenv("PASSWORD", "password123")
	require.Equal(t, expected, env.GetWMBPassword())
}

func TestGetWMBURL(t *testing.T) {
	expected := "http://example.com/wmb"
	os.Setenv("WMB_DIRECT_MESSAGE_URL", "http://example.com/wmb")
	require.Equal(t, expected, env.GetWMBURL())
}

func TestGetStateFilePath(t *testing.T) {
	expected := "/path/to/statefile"
	os.Setenv("STATEFILE", "/path/to/statefile")
	require.Equal(t, expected, env.GetStateFilePath())
}

func TestGetScheduledMinutes(t *testing.T) {
	expected := 5 * time.Minute
	require.Equal(t, expected, env.GetScheduledMinutes())
}

func TestGetFeedUrlsNotSet(t *testing.T) {
	os.Unsetenv("FEED_URLS")
	require.Equal(t, []string{""}, env.GetFeedUrls())
}

func TestGetFeedChannelsNotSet(t *testing.T) {
	os.Unsetenv("FEED_CHANNELS")
	require.Equal(t, []string{""}, env.GetFeedChannels())
}

func TestGetWMBPasswordNotSet(t *testing.T) {
	os.Unsetenv("PASSWORD")
	require.Equal(t, "", env.GetWMBPassword())
}

func TestGetWMBURLNotSet(t *testing.T) {
	os.Unsetenv("WMB_DIRECT_MESSAGE_URL")
	require.Equal(t, "", env.GetWMBURL())
}

func TestGetStateFilePathNotSet(t *testing.T) {
	os.Unsetenv("STATEFILE")
	require.Equal(t, "", env.GetStateFilePath())
}
