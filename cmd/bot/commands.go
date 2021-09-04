package main

import (
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/commands"
)

func registerCommands(tcb *bot.Bot) {
	tcb.Commands.RegisterCommand(commands.Ping(tcb))
	tcb.Commands.RegisterCommand(commands.Bot(tcb))
	tcb.Commands.RegisterCommand(commands.Game(tcb))
	tcb.Commands.RegisterCommand(commands.Title(tcb))
	tcb.Commands.RegisterCommand(commands.IsLive(tcb))
	tcb.Commands.RegisterCommand(commands.Events(tcb))
	tcb.Commands.RegisterCommand(commands.Help(tcb))
	tcb.Commands.RegisterCommand(commands.Subscribed(tcb))
	tcb.Commands.RegisterCommand(commands.NotifyMe(tcb))
	tcb.Commands.RegisterCommand(commands.RemoveMe(tcb))

	// Migration commands, admin use only
	// tcb.Commands.RegisterCommand(commands.MigrateChannels(tcb))
	// tcb.Commands.RegisterCommand(commands.MigrateSubscriptions(tcb))
}
