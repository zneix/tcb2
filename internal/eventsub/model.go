package eventsub

import (
	"encoding/json"
	"fmt"

	"github.com/nicklaw5/helix"
)

type eventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

type ChannelSubscription struct {
	Type      string
	Version   string
	ChannelID string
}

func (subscription *ChannelSubscription) String() string {
	return fmt.Sprintf("%s-%s@%s", subscription.Type, subscription.Version, subscription.ChannelID)
}

type EventSub struct {
	secret               string
	callbackURL          string
	onChannelUpdateEvent func(event helix.EventSubChannelUpdateEvent)
	onStreamOnlineEvent  func(event helix.EventSubStreamOnlineEvent)
	onStreamOfflineEvent func(event helix.EventSubStreamOfflineEvent)
}
