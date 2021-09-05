package commands

import (
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
)

func MigrateChannels(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "migrate_channels",
		Aliases:         []string{},
		Description:     "Migration command, admin use only",
		Usage:           "",
		IgnoreSelf:      false,
		CooldownChannel: 3 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			if msg.User.Name != "zneix" {
				return
			}

			// channel := tcb.Channels[msg.RoomID]
		},
	}
}
