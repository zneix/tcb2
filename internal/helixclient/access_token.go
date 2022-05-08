package helixclient

import (
	"log"
	"time"

	"github.com/nicklaw5/helix/v2"
)

// initAppAccessToken requests and sets app access token to the provided helix.Client
// and initializes a ticker running every 24 Hours which re-requests and sets app access token
func initAppAccessToken(client *helix.Client, tokenFetched chan struct{}) {
	response, err := client.RequestAppAccessToken([]string{})

	if err != nil {
		log.Fatalln("[Helix] Error requesting app access token:", err)
	}

	log.Printf("[Helix] Requested access token, status: %d, expires in: %d", response.StatusCode, response.Data.ExpiresIn)
	client.SetAppAccessToken(response.Data.AccessToken)
	close(tokenFetched)

	// Initialize the ticker
	ticker := time.NewTicker(24 * time.Hour)

	for range ticker.C {
		response, err := client.RequestAppAccessToken([]string{})
		if err != nil {
			log.Printf("[Helix] Failed to re-request app access token from ticker, status: %d", response.StatusCode)
			continue
		}
		log.Printf("[Helix] Re-requested access token from ticker, status: %d, expires in: %d", response.StatusCode, response.Data.ExpiresIn)

		client.SetAppAccessToken(response.Data.AccessToken)
	}
}
