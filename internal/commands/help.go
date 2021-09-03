package commands

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
)

func Help(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name: "help",
		Aliases: []string{
			"tcbhelp",
			"tcb_help",
			"titlechangebothelp",
			"titlechange_bothelp",
			"titlechangebot_help",
			"titlechange_bot_help",
		},
		Description:     "Posts a short list of commands or details about specified command",
		Usage:           "[command]",
		CooldownChannel: 100 * time.Millisecond,
		CooldownUser:    2 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]
			// Generic help
			if len(args) < 1 {
				cmdStrings := []string{}

				for _, cmd := range tcb.Commands.Commands {
					cmdStrings = append(cmdStrings, tcb.Commands.CommandString(cmd))
				}

				channel.Sendf("@%s, available commands: %s", msg.User.Name, strings.Join(cmdStrings, ", "))
				return
			}

			// Dynamic help
			cmd, exists := tcb.Commands.GetCommand(args[0])
			if !exists {
				channel.Sendf("@%s, provided command is either hidden or doesn't exist BrokeBack", msg.User.Name)
				return
			}
			description := strings.ReplaceAll(cmd.Description, "{prefix}", tcb.Commands.Prefix)
			channel.Sendf("@%s, %s (%s cooldown): %s", msg.User.Name, tcb.Commands.CommandString(cmd), cmd.CooldownUser, description)
		},
	}
}
