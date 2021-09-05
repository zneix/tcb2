package commands

import (
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/common"
)

func Ping(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "ping",
		Aliases:         []string{"tcbping"},
		Description:     "Pings the bot to see if it's online",
		Usage:           "",
		IgnoreSelf:      false,
		CooldownChannel: 1 * time.Second,
		CooldownUser:    2 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]
			channel.Sendf("@%s, reporting for duty MrDestructoid PowerUpR 🔔 %s", msg.User.Name, common.Version())
		},
	}
}
