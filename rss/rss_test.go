package rss_test

import (
	"testing"
	"time"

	"github.com/cfindlayisme/rss-wmb/db"
	"github.com/cfindlayisme/rss-wmb/rss"
	"github.com/stretchr/testify/require"
)

type mockFeedChecker struct {
	called   bool
	database *db.DB
	channels []string
	urls     []string
}

func (m *mockFeedChecker) CheckFeeds(database *db.DB, feedChannels []string, feedURLs []string) {
	m.called = true
	m.database = database
	m.channels = feedChannels
	m.urls = feedURLs
}

func TestScheduleFeeds(t *testing.T) {
	mock := &mockFeedChecker{}
	scheduler := &rss.Scheduler{FeedChecker: mock}

	channels := []string{"channel"}
	urls := []string{"url"}
	scheduler.ScheduleFeeds(nil, time.Millisecond, channels, urls)

	time.Sleep(10 * time.Millisecond)

	// Test that CheckFeeds was called
	require.True(t, mock.called)

	// Test that CheckFeeds was called with the correct arguments
	require.Nil(t, mock.database)
	require.Equal(t, channels, mock.channels)
	require.Equal(t, urls, mock.urls)
}

func TestScheduleFeedsNotCalledBeforeDuration(t *testing.T) {
	mock := &mockFeedChecker{}
	scheduler := &rss.Scheduler{FeedChecker: mock}

	scheduler.ScheduleFeeds(nil, time.Second, []string{"channel"}, []string{"url"})

	time.Sleep(500 * time.Millisecond)

	// Test that CheckFeeds was not called before the specified duration
	require.False(t, mock.called)
}

func TestScheduleFeedsCalledRepeatedly(t *testing.T) {
	mock := &mockFeedChecker{}
	scheduler := &rss.Scheduler{FeedChecker: mock}

	scheduler.ScheduleFeeds(nil, time.Millisecond, []string{"channel"}, []string{"url"})

	time.Sleep(10 * time.Millisecond)

	// Test that CheckFeeds was called
	require.True(t, mock.called)

	// Reset the called flag
	mock.called = false

	time.Sleep(10 * time.Millisecond)

	// Test that CheckFeeds was called again
	require.True(t, mock.called)
}
