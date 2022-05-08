package bot

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Run         func(msg twitch.PrivateMessage, args []string)

	CooldownChannel time.Duration
	CooldownUser    time.Duration

	// TODO: Perhaps cooldown logic should be stored in Bot / redis (?)
	LastExecutionChannel map[string]time.Time
	LastExecutionUser    map[string]time.Time
}

func (c *Command) String() string {
	str := c.Name
	if c.Usage != "" {
		str += " " + c.Usage
	}
	return str
}

type CommandController struct {
	Commands map[string]*Command
	aliases  map[string]string
	Prefix   string
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
