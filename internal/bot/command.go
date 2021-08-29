package bot

import (
	"strings"
	"time"
)

func (c *CommandController) GetCommand(alias string) (*Command, bool) {
	name, ok := c.aliases[strings.ToLower(alias)]
	if !ok {
		return nil, false
	}

	cmd, ok := c.commands[name]
	return cmd, ok
}

func (c *CommandController) RegisterCommand(cmd *Command) {
	cmd.LastExecutionChannel = make(map[string]time.Time)
	cmd.LastExecutionUser = make(map[string]time.Time)

	c.commands[cmd.Name] = cmd

	c.aliases[cmd.Name] = cmd.Name
	for _, alias := range cmd.Aliases {
		c.aliases[alias] = cmd.Name
	}
}

func NewCommandController() *CommandController {
	return &CommandController{
		commands: make(map[string]*Command),
		aliases:  make(map[string]string),
	}
}
