package supinicapi

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func New(supinicAPIKey string) *Client {
	return &Client{
		apiKey:     supinicAPIKey,
		httpClient: &http.Client{},
	}
}

// requestAliveStatus check documentation at https://supinic.com/api/#api-Bot_Program-UpdateBotActivity
func (c *Client) requestAliveStatus() {
	req, err := http.NewRequest("PUT", "https://supinic.com/api/bot-program/bot/active", nil)
	if err != nil {
		log.Printf("[SupinicAPI] Error creating API reqeust: %s\n", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.apiKey))

	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[SupinicAPI] Failed to update alive status: %s\n", err)
		return
	}
	defer res.Body.Close()

	log.Println(res)
	log.Printf("[SupinicAPI] Successfully updated alive status")
}

func (c *Client) UpdateAliveStatus() {
	c.requestAliveStatus()

	// Make the API request every 15 minutes
	ticker := time.NewTicker(15 * time.Minute)

	for range ticker.C {
		c.requestAliveStatus()
	}
}
