package main

import "github.com/zneix/tcb2/internal/bot"

// TODO: Add function creating channels

// joinChannels performs startup actions for all the channels that are already loaded
func joinChannels(tcb *bot.Bot) {
	for ID, channel := range tcb.Channels {
		// Set the ID in map translating login names back to IDs
		tcb.Logins[channel.Login] = ID

		// Start message queue routine
		// TODO: this should be moved to where a Channel is created
		//go channel.StartMessageQueue(tcb.TwitchIRC)

		// JOIN the channel
		tcb.TwitchIRC.Join(channel.Login)
	}
}
