package helixclient

import (
	"errors"

	"github.com/nicklaw5/helix"
	"github.com/zneix/tcb2/internal/config"
)

// New returns a helix.Client that has requested an AppAccessToken and will keep it refreshed every 24h
func New(cfg config.TCBConfig) (*helix.Client, error) {
	if cfg.TwitchClientID == "" {
		return nil, errors.New("Twitch Client ID is missing, can't make Helix requests")
	}
	if cfg.TwitchClientSecret == "" {
		return nil, errors.New("Twitch Client Secret is missing, can't make Helix requests")
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.TwitchClientID,
		ClientSecret: cfg.TwitchClientSecret,
	})
	if err != nil {
		return nil, err
	}

	// Initialize methods responsible for refreshing access token
	waitForFirstAppAccessToken := make(chan struct{})
	go initAppAccessToken(client, waitForFirstAppAccessToken)
	<-waitForFirstAppAccessToken

	return client, nil
}
