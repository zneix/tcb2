package bot

import (
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/nicklaw5/helix"
	"github.com/zneix/tcb2/internal/eventsub"
	"github.com/zneix/tcb2/internal/mongo"
)

// types

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

type Channel struct {
	ID    string      `bson:"id"`
	Login string      `bson:"login"`
	Mode  ChannelMode `bson:"mode"`

	LastMsg      string
	QueueChannel chan *QueueMessage
}

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Run         func(msg twitch.PrivateMessage, args []string)

	CooldownChannel time.Duration
	CooldownUser    time.Duration

	// TODO: Perhaps cooldown logic should be stored in Bot / redis (?)
	LastExecutionChannel map[string]time.Time
	LastExecutionUser    map[string]time.Time

	// TODO properties:
	// Usage string
}

type CommandController struct {
	commands map[string]*Command
	aliases  map[string]string
}

type QueueMessage struct {
	Message string
}

// enums

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
	//ChannelModeVIP

	// ChannelModeEnumBoundary used mark the end of enumeration
	ChannelModeEnumBoundary
)

// MessageRatelimit the minimum time.Duration that must pass between sending messages in the Channel
func (mode ChannelMode) MessageRatelimit() time.Duration {
	if mode == ChannelModeModerator {
		// 1200ms is minimum, but 1650ms prevents exceeding global limits
		return 1650 * time.Millisecond
	}
	return 100 * time.Millisecond
}

func (mode ChannelMode) MessageLengthMax() int {
	if mode == ChannelModeModerator {
		// Leaving 2 characters for the magic character
		return 498
	}
	// TODO: Investigate the actual limit for "pleb" modes (?)
	// mm2pl: maybe it's something like max of count(CHAR) / len(msg) for each unique character used in a message
	// mm2pl: or maybe it's some kind of GOW average
	// mm2pl: max((msg.count(ch) / len(msg) for ch in set(msg))) seems like a good approximation
	// For now I'm lazy and just gonna hardcode some reasonable value in here
	return 468
}
