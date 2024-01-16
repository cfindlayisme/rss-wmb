package wmb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cfindlayisme/rss-wmb/env"
	"github.com/cfindlayisme/rss-wmb/model"
	"github.com/mmcdole/gofeed"
)

func SendDirectedRSSMessage(item *gofeed.Item, feedChannels []string, n int) {
	log.Printf("Title: %s\n", item.Title)
	log.Println("--------------------")

	// Create a new WebhookMessage
	webhookMessage := model.WebhookMessage{
		Message:  "Title: " + item.Title + " Link: " + item.Link,
		Password: env.GetWMBPassword(),
	}

	webhookDirectedMessage := model.DirectedWebhookMessage{
		IncomingMessage: webhookMessage,
		Target:          feedChannels[n],
	}

	jsonData, err := json.Marshal(webhookDirectedMessage)
	if err != nil {
		log.Fatalf("Error marshalling webhookDirectedMessage: %v", err)
	}

	log.Printf("JSON Data: %s\n", jsonData)

	// Send a POST request to the webhook URL
	resp, err := http.Post(env.GetWMBURL(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error sending webhook: %v", err)
	}
	// Read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the response body
	log.Printf("Response: %s\n", body)
	defer resp.Body.Close()
}
