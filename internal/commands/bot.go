package commands

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
)

func Bot(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "bot",
		Aliases:         []string{"tcb", "about", "titlechangebot", "titlechange_bot"},
		Description:     "Returns basic information about the bot",
		Usage:           "",
		CooldownChannel: 3 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]
			channel.Send(fmt.Sprintf("I am a bot created by zneix. I can notify you when the channel goes live or the title changes. Try %shelp for a list of commands. pajaDank", msg.User.Name))
		},
	}
}
