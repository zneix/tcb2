package commands

import (
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/zneix/tcb2/internal/bot"
)

func Title(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "title",
		Aliases:         []string{"currenttitle"},
		Description:     "Returns current title",
		Usage:           "",
		CooldownChannel: 3 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]
			channel.Sendf("@%s, current title: %s", msg.User.Name, channel.CurrentTitle)
		},
	}
}
