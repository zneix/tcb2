package supinicapi

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Client represents a Supinic's API wrapper
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// New creates a new Client instance, requires a
func New(supinicAPIKey string) *Client {
	return &Client{
		apiKey:     supinicAPIKey,
		httpClient: &http.Client{},
	}
}

// requestAliveStatus check documentation at https://supinic.com/api/#api-Bot_Program-UpdateBotActivity
func (c *Client) requestAliveStatus() {
	req, err := http.NewRequest("PUT", "https://supinic.com/api/bot-program/bot/active", http.NoBody)
	if err != nil {
		log.Printf("[SupinicAPI] Error creating API request: %s\n", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.apiKey))

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[SupinicAPI] Failed to update alive status: %s\n", err)
		return
	}
	defer res.Body.Close()

	log.Printf("[SupinicAPI] Pinged alive endpoint, status: %d\n", res.StatusCode)
}

// UpdateAliveStatus starts routine updating alive status on Supinic's API right away and every 15 minutes
func (c *Client) UpdateAliveStatus() {
	if c.apiKey == "" {
		log.Println("[SupinicAPI] API key is empty, won't make API requests")
		return
	}

	c.requestAliveStatus()

	// Make the API request every 15 minutes
	ticker := time.NewTicker(15 * time.Minute)

	for range ticker.C {
		c.requestAliveStatus()
	}
}
