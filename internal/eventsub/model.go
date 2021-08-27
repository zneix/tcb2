package eventsub

import (
	"encoding/json"

	"github.com/nicklaw5/helix"
)

type eventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

type EventSub struct {
	secret               string
	onChannelUpdateEvent func(event helix.EventSubChannelUpdateEvent)
	onStreamOnlineEvent  func(event helix.EventSubStreamOnlineEvent)
	onStreamOfflineEvent func(event helix.EventSubStreamOfflineEvent)
}
