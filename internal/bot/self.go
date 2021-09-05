package bot

import (
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/config"
)

func (s *Self) JoinChannel(twitchIRC *twitch.Client) {
	twitchIRC.Join(s.Channel.Login)
}

func NewSelf(cfg *config.TCBConfig, twitchIRC *twitch.Client) *Self {
	channel := &Channel{
		ID:                 "", // We don't really need the ID for anything, but maybe it should be defined just in case
		Login:              cfg.TwitchLogin,
		DisabledCommands:   make([]string, 0),
		Events:             make(map[SubEventType]string),
		PajbotAPI:          nil,
		MessageLengthLimit: 0,
		WhisperCommands:    false,
		EventsOnlyOffline:  false,
		Mode:               ChannelModeModerator,
		QueueChannel:       make(chan *QueueMessage),
	}

	go channel.StartMessageQueue(twitchIRC)

	return &Self{
		Login: cfg.TwitchLogin,
		// OAuth:   cfg.TwitchOAuth,
		Channel: channel,
	}
}
