package main

import (
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/commands"
)

func registerCommands(tcb *bot.Bot) {
	tcb.Commands.RegisterCommand(commands.Ping(tcb))
}
