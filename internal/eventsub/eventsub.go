package eventsub

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/nicklaw5/helix"
	"github.com/zneix/tcb2/internal/api"
	"github.com/zneix/tcb2/internal/config"
)

func (esub *EventSub) CreateChannelSubscription(helixClient *helix.Client, subscription *ChannelSubscription) error {
	resp, err := helixClient.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    subscription.Type,
		Version: subscription.Version,
		Condition: helix.EventSubCondition{
			BroadcasterUserID: subscription.ChannelID,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: esub.callbackURL,
			Secret:   esub.secret,
		},
	})
	if err != nil {
		return err
	}

	log.Printf("[EventSub] Create subscription response for %s: %# v\n", subscription, resp.Data)

	// TODO: Properly handle pending status
	//subscriptionsPending = append(subscriptionsPending, sub.ID)

	return nil
}

func (esub *EventSub) handleIncomingNotification(notification *eventSubNotification) {
	switch notification.Subscription.Type {

	case helix.EventSubTypeChannelUpdate:
		// channel.update
		var event helix.EventSubChannelUpdateEvent
		err := json.Unmarshal(notification.Event, &event)
		if err != nil {
			log.Printf("[EventSub] Failed to unmarshal notification event: %s, data: %s\n", err, string(notification.Event))
			return
		}
		if esub.onChannelUpdateEvent != nil {
			esub.onChannelUpdateEvent(event)
		}

	case helix.EventSubTypeStreamOnline:
		// stream.online
		var event helix.EventSubStreamOnlineEvent
		err := json.Unmarshal(notification.Event, &event)
		if err != nil {
			log.Printf("[EventSub] Failed to unmarshal notification event: %s, data: %s\n", err, string(notification.Event))
			return
		}
		if esub.onStreamOnlineEvent != nil {
			esub.onStreamOnlineEvent(event)
		}

	case helix.EventSubTypeStreamOffline:
		// stream.offline
		var event helix.EventSubStreamOfflineEvent
		err := json.Unmarshal(notification.Event, &event)
		if err != nil {
			log.Printf("[EventSub] Failed to unmarshal notification event: %s, data: %s\n", err, string(notification.Event))
			return
		}
		if esub.onStreamOfflineEvent != nil {
			esub.onStreamOfflineEvent(event)
		}

	default:
		log.Printf("[EventSub] Received unhandled notification: %# v\n", notification)
	}
}

// OnChannelUpdateEvent attach callback to channel.update event
func (esub *EventSub) OnChannelUpdateEvent(callback func(event helix.EventSubChannelUpdateEvent)) {
	esub.onChannelUpdateEvent = callback
}

// OnStreamOnlineEvent attach callback to stream.online event
func (esub *EventSub) OnStreamOnlineEvent(callback func(event helix.EventSubStreamOnlineEvent)) {
	esub.onStreamOnlineEvent = callback
}

// OnStreamOfflineEvent attach callback to stream.offline event
func (esub *EventSub) OnStreamOfflineEvent(callback func(event helix.EventSubStreamOfflineEvent)) {
	esub.onStreamOfflineEvent = callback
}

func New(cfg *config.TCBConfig, apiServer *api.Server) *EventSub {
	eventsub := &EventSub{
		secret:      cfg.TwitchEventSubSecret,
		callbackURL: strings.TrimSuffix(apiServer.BaseURL, "/") + "/eventsub/callback",
	}

	eventsub.registerAPIRoutes(apiServer)

	return eventsub
}
