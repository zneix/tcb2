package main

import (
	"log"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/nicklaw5/helix"
	"github.com/zneix/tcb2/internal/bot"
)

func initializeEvents(tcb *bot.Bot) {
	// Twitch IRC events

	// Authenticated with IRC
	tcb.TwitchIRC.OnConnect(func() {
		log.Println("[TwitchIRC] connected")
		joinChannels(tcb)
	})

	// PRIVMSG
	tcb.TwitchIRC.OnPrivateMessage(func(message twitch.PrivateMessage) {
		// TODO: Add PRIVMSG handling
		//channel := tcb.Channels[message.RoomID]
	})

	// USERSTATE
	tcb.TwitchIRC.OnUserStateMessage(func(message twitch.UserStateMessage) {
		channelID, ok := tcb.Logins[message.Channel]
		if !ok {
			// tcb.Logins map didn't have current channel's ID
			// Note: this should realistically never occur though, but early exit to prevent panic
			return
		}

		channel := tcb.Channels[channelID]

		// Check if Channel.Mode changed by comparing bot's state
		newMode := bot.ChannelModeNormal

		// Bot will always have elevated permissions in its own chat, saving some time with the early-out
		if channel.Login == tcb.Self.Login {
			return
		}

		// First check user-type
		userType, ok := message.Tags["user-type"]
		if !ok {
			log.Println("[USERSTATE] user-type tag was not found in the IRC message, either no capabilities or Twitch removed this tag xd")
		} else if userType == "mod" {
			newMode = bot.ChannelModeModerator
		} else {
			// Since user-type does not care about VIP status, we need to check badges
			for key := range message.User.Badges {
				if key == "vip" || key == "moderator" {
					newMode = bot.ChannelModeModerator
					break
				}
			}

		}

		// Update ChannelMode in the current channel if it differs
		if newMode != channel.Mode {
			channel.ChangeMode(tcb.Mongo, newMode)
		}

	})

	// Twitch EventSub events

	// channel.update
	tcb.EventSub.OnChannelUpdateEvent(func(event helix.EventSubChannelUpdateEvent) {
		// TODO: Handle received event
	})

	// stream.online
	tcb.EventSub.OnStreamOnlineEvent(func(event helix.EventSubStreamOnlineEvent) {
		// TODO: Handle received event
	})

	// stream.offline
	tcb.EventSub.OnStreamOfflineEvent(func(event helix.EventSubStreamOfflineEvent) {
		// TODO: Handle received event
	})
}
