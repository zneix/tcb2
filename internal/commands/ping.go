package commands

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
)

func Ping(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:        "ping",
		Aliases:     []string{"tcbping"},
		Description: "",
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]
			channel.Send(fmt.Sprintf("@%s, reporting for duty NaM 7", msg.User.Name))
		},
	}
}
