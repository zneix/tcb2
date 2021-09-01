package bot

import (
	"strings"
	"time"
)

func (c *CommandController) CommandString(cmd *Command) string {
	str := c.Prefix + cmd.Name
	if cmd.Usage != "" {
		str += " " + cmd.Usage
	}
	return str
}

func (c *CommandController) GetCommand(alias string) (*Command, bool) {
	name, ok := c.aliases[strings.ToLower(alias)]
	if !ok {
		return nil, false
	}

	cmd, ok := c.Commands[name]
	return cmd, ok
}

func (c *CommandController) RegisterCommand(cmd *Command) {
	cmd.LastExecutionChannel = make(map[string]time.Time)
	cmd.LastExecutionUser = make(map[string]time.Time)

	c.Commands[cmd.Name] = cmd

	c.aliases[cmd.Name] = cmd.Name
	for _, alias := range cmd.Aliases {
		c.aliases[alias] = cmd.Name
	}
}

func NewCommandController(prefix string) *CommandController {
	return &CommandController{
		Commands: make(map[string]*Command),
		aliases:  make(map[string]string),
		Prefix:   prefix,
	}
}
