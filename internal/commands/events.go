package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
)

func Events(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "events",
		Aliases:         []string{"tcbevents"},
		Description:     "Shows available events you can subscribe to with a brief description",
		Usage:           "",
		CooldownChannel: 3 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]

			eventStrings := []string{}
			for i, desc := range bot.SubEventDescriptions {
				eventStrings = append(eventStrings, fmt.Sprintf("%s (%s)", bot.SubEventType(i), desc))
			}
			// TODO: Export default command prefix to somewhere, e.g. config
			channel.Send(fmt.Sprintf("@%s, available events: %s. Type \"%snotifyme <event> [optional value]\" to subscribe to an event!", msg.User.Name, strings.Join(eventStrings, ", "), tcb.Commands.Prefix))
		},
	}
}
