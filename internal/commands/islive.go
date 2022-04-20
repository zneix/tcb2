package commands

import (
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/zneix/tcb2/internal/bot"
)

func IsLive(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "islive",
		Aliases:         []string{"tcbislive"},
		Description:     "Shows you if the channel is live or not",
		Usage:           "",
		CooldownChannel: 3 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]

			var isLive = channel.IsLive
			if len(args) >= 1 {
				targetID, ok := tcb.Logins[args[0]]
				if !ok {
					channel.Send("I don't track the status of the target channel")
					return
				}
				isLive = tcb.Channels[targetID].IsLive
			}
			var message string
			if isLive {
				message = "Target channel is live KKona GuitarTime"
			} else {
				message = "Target channel is offline FeelsBadMan TeaTime"
			}
			channel.Send(message)
		},
	}
}
