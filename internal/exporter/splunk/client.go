package splunk

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Client struct {
	URL   string
	Token string
	// Buffer ...
}

func NewClient(url, token string) *Client {
	return &Client{
		URL:   url,
		Token: token,
	}
}

// SendEvent sends a single raw event to Splunk HEC.
func (c *Client) SendEvent(event interface{}) error {
	payload := map[string]interface{}{
		"event": event,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Splunk "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Splunk HEC error: %s", resp.Status)
	}

	return nil
}
