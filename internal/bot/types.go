package bot

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"github.com/zneix/tcb2/internal/eventsub"
	"github.com/zneix/tcb2/internal/mongo"
)

// types
// TODO: Restructure this and split types to their separate types

// Self contains properties related to bot's user account
type Self struct {
	Login string
	OAuth string
}

type Bot struct {
	TwitchIRC *twitch.Client
	Mongo     *mongo.Connection
	Helix     *helix.Client
	EventSub  *eventsub.EventSub

	Logins   map[string]string
	Channels map[string]*Channel
	Commands *CommandController

	Self      *Self
	StartTime time.Time
}

type PajbotAPI struct {
	Mode   PajbotAPIMode `bson:"mode"`
	Domain string        `bson:"domain"`
}

// SubEventMOTD if present for a channel with the corresponding ChannelID, should be posted right after announcing channel going live
// It could be useful to remind streamer to tweet or announce going live on Discord
type SubEventMOTD struct {
	ChannelID string `bson:"channel_id"`
	Message   string `bson:"message"`
}

type SubEventSubscription struct {
	UserLogin string       `bson:"user_login"`
	UserID    string       `bson:"user_id"`
	Event     SubEventType `bson:"event"`
	Value     string       `bson:"value"`
}

type SubEventMessage struct {
	Bot       *Bot
	ChannelID string
	Type      SubEventType
}

type QueueMessage struct {
	Message string
}

// enums

//
// ChannelMode indicates the bot's state in a Channel
type ChannelMode int

const (
	// ChannelModeNormal default ChannelMode with regular ratelimits
	ChannelModeNormal ChannelMode = iota
	// ChannelModeInactive bot has been disabled in that chat
	ChannelModeInactive
	// ChannelModeModerator bot has elevated ratelimits with moderation permissions
	ChannelModeModerator
	// ChannelModeVIP bot has elevated ratelimits without moderation permissions
	// Note: we don't need this, but maybe it can be useful in the future
	// ChannelModeVIP

	// ChannelModeEnumBoundary marks the end of enumeration
	ChannelModeEnumBoundary
)

// MessageRatelimit the minimum time.Duration that must pass between sending messages in the Channel
func (mode ChannelMode) MessageRatelimit() time.Duration {
	if mode == ChannelModeModerator {
		return 100 * time.Millisecond
	}
	// 1200ms is minimum, but 1650ms prevents exceeding global limits
	return 1650 * time.Millisecond
}

//
// PajbotAPIMode indicates bot's behavior regarding banphrase checks in channels that have pajbot API configured
type PajbotAPIMode int

const (
	// PajbotAPIModeDisabled even if the Domain link is set, API will be totally ignored
	PajbotAPIModeDisabled PajbotAPIMode = iota
	// PajbotAPIModeEnabled will attempt to sanitize potentially harmful message content
	PajbotAPIModeEnabled
)

//
// SubEventType defines event to which users can subscribe
type SubEventType int

const (
	// SubEventTypeGame game (category) has been updated
	// Received in EventSub's "channel.update"
	SubEventTypeGame SubEventType = iota
	// SubEventTypeTitle title has been updated
	// Received in EventSub's "channel.update"
	SubEventTypeTitle
	// SubEventTypeLive channel has gone live
	// Received in EventSub's "stream.online"
	SubEventTypeLive
	// EventLive channel has gone offline
	// Received in EventSub's "stream.offline"
	SubEventTypeOffline
	// SubEventTypePartnered broadcaster has become partnered
	// This is deprecated and is kept in case legacy subscriptions with this value are found
	SubEventTypePartnered

	// SubEventTypeInvalid represents an invalid event type that was passed to ParseChannelEvent
	SubEventTypeInvalid
)

func (e SubEventType) String() string {
	switch e {
	case SubEventTypeGame:
		return "game"
	case SubEventTypeTitle:
		return "title"
	case SubEventTypeLive:
		return "live"
	case SubEventTypeOffline:
		return "offline"
	case SubEventTypePartnered:
		return "partnered"
	default:
		return fmt.Sprintf("invalid(%d)", e)
	}
}
