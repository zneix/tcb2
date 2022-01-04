package main

import (
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/zneix/tcb2/internal/bot"
)

func registerEvents(tcb *bot.Bot) {
	// Twitch IRC events

	// Authenticated with IRC
	tcb.TwitchIRC.OnConnect(func() {
		log.Println("[TwitchIRC] connected")
		joinChannels(tcb)
	})

	// PRIVMSG
	tcb.TwitchIRC.OnPrivateMessage(func(message twitch.PrivateMessage) {
		pajbotAlert(tcb, &message)

		// Ignore non-commands
		if !strings.HasPrefix(message.Message, tcb.Commands.Prefix) {
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

		userType, ok := message.Tags["user-type"]
		switch {
		case !ok:
			log.Println("[USERSTATE] user-type tag was not found in the IRC message, either no capabilities or Twitch removed this tag xd")

		case userType == "mod":
			newMode = bot.ChannelModeModerator

		default:
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
			err := channel.ChangeMode(tcb.Mongo, newMode)
			if err != nil {
				log.Printf("Failed to change mode in %s: %s\n", channel, err)
			}
		}
	})

	// NOTICE
	tcb.TwitchIRC.OnNoticeMessage(func(message twitch.NoticeMessage) {
		channelID, ok := tcb.Logins[message.Channel]
		if !ok {
			// tcb.Logins map didn't have current channel's ID
			// Note: this should realistically never occur though, but early exit to prevent panic
			return
		}
		channel := tcb.Channels[channelID]

		log.Printf("[TwitchIRC:NOTICE] %s in %s\n", message.MsgID, channel)

		switch message.MsgID {
		case "msg_banned", "msg_channel_suspended":
			err := channel.ChangeMode(tcb.Mongo, bot.ChannelModeInactive)
			if err != nil {
				log.Printf("Failed to change mode in %s: %s\n", channel, err)
			}
		default:
		}
	})

	// Twitch EventSub events

	// channel.update
	tcb.EventSub.OnChannelUpdateEvent(func(event helix.EventSubChannelUpdateEvent) {
		log.Printf("[EventSub:channel.update] %# v\n", event)
		channel := tcb.Channels[event.BroadcasterUserID]

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
		channel := tcb.Channels[event.BroadcasterUserID]

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
		channel := tcb.Channels[event.BroadcasterUserID]

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

// Little fun module responding to pajbot alerts
func pajbotAlert(tcb *bot.Bot, msg *twitch.PrivateMessage) {
	if msg.RoomID == "11148817" && (msg.User.ID == "82008718" || msg.User.ID == "99631238") {
		if msg.Action && strings.HasPrefix(msg.Message, "pajaS 🚨 ALERT") {
			log.Printf("[pajbotAlert] triggered by %s(%s)", msg.User.Name, msg.User.ID)

			channel := tcb.Channels[msg.RoomID]
			channel.Send(".me pajaDinkDonk 🚨 PINGS")
		}
	}
}
