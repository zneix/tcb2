package main

import (
	"log"

	"github.com/gempir/go-twitch-irc/v2"
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
		if channel.Login == tcb.Self.Login {
			// Bot will always have elevated permissions in its own chat, saving some time with the early-out
			return
		}

		// Check if bot's state in the Channel has changed
		newMode := bot.ChannelModeNormal
		for key := range message.User.Badges {
			if key == "moderator" {
				newMode = bot.ChannelModeModerator
				break
			}
		}

		// Update ChannelMode in the current channel if it differs
		if newMode != channel.Mode {
			// TODO: Acknowledge mode change
			channel.Mode = newMode
		}

	})
}
