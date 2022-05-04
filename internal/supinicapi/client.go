package supinicapi

import (
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
		log.Println("[SupinicAPI] Error creating API request:", err)
		return
	}
	req.Header.Set("Authorization", "Basic "+c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("[SupinicAPI] Failed to update alive status:", err)
		return
	}
	defer res.Body.Close()

	log.Println("[SupinicAPI] Pinged alive endpoint, status:", res.StatusCode)
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
