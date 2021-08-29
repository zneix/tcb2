package commands

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
)

func Game(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "game",
		Aliases:         []string{"currentgame"},
		Description:     "Returns current game",
		CooldownChannel: 3 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]
			channel.Send(fmt.Sprintf("@%s, current game: %s", msg.User.Name, channel.CurrentGame))
		},
	}
}
