package main

import (
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
	"github.com/zneix/tcb2/internal/bot"
)

// handlerOnNoticeMessage logic for NOTICE twitch IRC messages
// It's taken out of registerEvents since it's used for both read and write conns
func handlerOnNoticeMessage(tcb *bot.Bot, message *twitch.NoticeMessage, connType string) {
	channelID, ok := tcb.Logins[message.Channel]
	if !ok {
		// tcb.Logins map didn't have current channel's ID
		// Note: this should realistically never occur though, but early exit to prevent panic
		return
	}
	channel := tcb.Channels[channelID]

	log.Printf("[TwitchIRC:%s] NOTICE %s in %s: %s\n", connType, message.MsgID, channel, message.Message)

	switch message.MsgID {
	case "msg_banned", "msg_channel_suspended":
		err := channel.ChangeMode(tcb.Mongo, bot.ChannelModeInactive)
		if err != nil {
			log.Printf("Failed to change mode in %s: %s\n", channel, err)
		}
	default:
	}
}

func registerEvents(tcb *bot.Bot) {
	// Twitch IRC events

	// Authenticated with IRC
	tcb.TwitchRead.OnConnect(func() {
		log.Println("[TwitchIRC:read] connected, joining channels")
		joinChannels(tcb)
	})
	tcb.TwitchWrite.OnConnect(func() {
		log.Println("[TwitchIRC:write] connected")
	})

	tcb.TwitchRead.OnReconnectMessage(func(message twitch.ReconnectMessage) {
		log.Println("[TwitchIRC:read] received RECONNECT:", message.Raw)
	})

	// PRIVMSG
	tcb.TwitchRead.OnPrivateMessage(func(message twitch.PrivateMessage) {
		// Early out in case message does not start with command prefix - meaning it's not a command
		if !strings.HasPrefix(message.Message, tcb.Commands.Prefix) {
			// Handle non-commands
			return
		}

		// Parse command name and arguments
		args := strings.Fields(message.Message)
		commandName := strings.ToLower(args[0][utf8.RuneCountInString(tcb.Commands.Prefix):])
		args = args[1:]

		// Try to find the command by its name and/or aliases
		command, exists := tcb.Commands.GetCommand(commandName)
		if !exists {
			return
		}

		// Skip command execution if it's disabled in the target channel
		channel := tcb.Channels[message.RoomID]
		for _, disabledCmdName := range channel.DisabledCommands {
			if commandName == disabledCmdName {
				return
			}
		}

		// TODO: [Permissions] Check if user is allowed to execute the command

		// Check if channel or user is on cooldown
		if time.Since(command.LastExecutionChannel[message.RoomID]) < command.CooldownChannel || time.Since(command.LastExecutionUser[message.User.ID]) < command.CooldownUser {
			return
		}

		// Execute the command
		go command.Run(message, args)

		// Apply cooldown if user's permissions don't allow to skip it
		// TODO: [Permissions] Don't apply user cooldowns to users that are allowed to skip it
		command.LastExecutionChannel[message.RoomID] = time.Now()
		command.LastExecutionUser[message.User.ID] = time.Now()
	})

	// USERSTATE
	// These will be triggered whenever a message is written to a channel - so react to those on write connection
	// They might also be received upon JOINing on authed connection, however we don't do that
	tcb.TwitchWrite.OnUserStateMessage(func(message twitch.UserStateMessage) {
		channelID, ok := tcb.Logins[message.Channel]
		if !ok {
			// tcb.Logins map didn't have current channel's ID
			// Note: this should realistically never occur though, but early exit to prevent panic
			return
		}

		channel := tcb.Channels[channelID]

		// Bot will always have elevated permissions in its own chat, saving some time with the early-out
		if channel.Login == tcb.Self.Login {
			return
		}

		// Check if Channel.Mode changed to see if we now have privileged write limits - by being either a moderator or vip
		newMode := bot.ChannelModeNormal

		if message.User.IsMod || message.User.IsVip {
			newMode = bot.ChannelModeModerator
		}

		// Update ChannelMode in the current channel if it differs
		if newMode != channel.Mode {
			err := channel.ChangeMode(tcb.Mongo, newMode)
			if err != nil {
				log.Printf("Failed to change mode in %s: %s\n", channel, err)
			}
		}
	})

	// NOTICE
	// This might be relevant for both read and write connections:
	// on Read: a channel might be suspended
	// on Write: the bot user might be banned from channel it attempts to send a message in
	tcb.TwitchRead.OnNoticeMessage(func(message twitch.NoticeMessage) {
		handlerOnNoticeMessage(tcb, &message, "read")
	})
	tcb.TwitchWrite.OnNoticeMessage(func(message twitch.NoticeMessage) {
		handlerOnNoticeMessage(tcb, &message, "write")
	})

	// Twitch EventSub events

	// channel.update
	tcb.EventSub.OnChannelUpdateEvent(func(event helix.EventSubChannelUpdateEvent) {
		log.Printf("[EventSub:channel.update] %# v\n", event)
		channel, ok := tcb.Channels[event.BroadcasterUserID]
		if !ok {
			return
		}

		// Announce game change
		if event.CategoryName != channel.CurrentGame {
			channel.CurrentGame = event.CategoryName
			go subEventTrigger(&bot.SubEventMessage{
				Bot:       tcb,
				ChannelID: event.BroadcasterUserID,
				Type:      bot.SubEventTypeGame,
			})
		}
		// Announce title change
		if event.Title != channel.CurrentTitle {
			channel.CurrentTitle = event.Title
			go subEventTrigger(&bot.SubEventMessage{
				Bot:       tcb,
				ChannelID: event.BroadcasterUserID,
				Type:      bot.SubEventTypeTitle,
			})
		}
	})

	// stream.online
	tcb.EventSub.OnStreamOnlineEvent(func(event helix.EventSubStreamOnlineEvent) {
		log.Printf("[EventSub:stream.online] %# v\n", event)
		channel, ok := tcb.Channels[event.BroadcasterUserID]
		if !ok {
			return
		}

		if channel.IsLive {
			log.Printf("[EventSub] Received stream.online, but %s seems to be already live!", channel)
			return
		}

		channel.IsLive = true
		// Announce channel going live
		go subEventTrigger(&bot.SubEventMessage{
			Bot:       tcb,
			ChannelID: event.BroadcasterUserID,
			Type:      bot.SubEventTypeLive,
		})
	})

	// stream.offline
	tcb.EventSub.OnStreamOfflineEvent(func(event helix.EventSubStreamOfflineEvent) {
		log.Printf("[EventSub:stream.offline] %# v\n", event)
		channel, ok := tcb.Channels[event.BroadcasterUserID]
		if !ok {
			return
		}

		if !channel.IsLive {
			log.Printf("[EventSub] Received stream.offline, but %s seems to be already offline!", channel)
			return
		}

		channel.IsLive = false
		// Announce channel going offline
		go subEventTrigger(&bot.SubEventMessage{
			Bot:       tcb,
			ChannelID: event.BroadcasterUserID,
			Type:      bot.SubEventTypeOffline,
		})
	})
}
