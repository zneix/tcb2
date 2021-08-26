package main

import (
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/config"
)

const (
	version = "2.0-alpha"
)

func main() {
	log.Printf("Starting titlechange_bot v%s", version)

	cfg := config.New()

	self := &bot.Self{
		Login: cfg.TwitchLogin,
		OAuth: cfg.TwitchOAuth,
	}

	tcb := &bot.Bot{
		TwitchIRC: twitch.NewClient(cfg.TwitchLogin, "oauth:"+cfg.TwitchOAuth),
		Self:      self,
		Logins:    make(map[string]string),
		Channels:  make(map[string]*bot.Channel),
		Commands:  make(map[string]*bot.Command),
		StartTime: time.Now(),
	}

	// init actions that require bot.Bot object initialized already
	initializeEvents(tcb)

	err := tcb.TwitchIRC.Connect()
	if err != nil {
		log.Fatalln(err)
	}
}
