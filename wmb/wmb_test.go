package wmb_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cfindlayisme/rss-wmb/wmb"
	"github.com/mmcdole/gofeed"
)

func TestSendDirectedRSSMessage(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test method and path
		if req.URL.String() != "/" || req.Method != http.MethodPost {
			t.Errorf("Expected POST method; got %s", req.Method)
		}

		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))

	defer server.Close()

	// Test success case with a good url
	item := &gofeed.Item{Title: "Test Title", Link: "Test Link"}
	feedChannels := []string{"Test Channel"}
	err := wmb.SendDirectedRSSMessage(server.URL, item, feedChannels, 0)
	if err != nil {
		t.Errorf("Expected no error; got %v", err)
	}

	// Test failure case with an invalid url
	err = wmb.SendDirectedRSSMessage("http://invalid-url", item, feedChannels, 0)
	if err == nil {
		t.Error("Expected error; got none")
	}
}
