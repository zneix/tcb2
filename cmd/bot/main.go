package main

import (
	"context"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/config"
	"github.com/zneix/tcb2/internal/helixclient"
	"github.com/zneix/tcb2/internal/mongo"
)

const (
	VERSION = "2.0-alpha"
)

func main() {
	log.Printf("Starting titlechange_bot v%s", VERSION)

	cfg := config.New()
	ctx := context.Background()

	mongoConnection := mongo.NewMongoConnection(cfg, ctx)
	mongoConnection.Connect()

	twitchIRC := twitch.NewClient(cfg.TwitchLogin, "oauth:"+cfg.TwitchOAuth)

	helixClient, err := helixclient.New(cfg)
	if err != nil {
		log.Fatalf("[Helix] Error while initializing client: %s\n", err)
	}

	self := &bot.Self{
		Login: cfg.TwitchLogin,
		OAuth: cfg.TwitchOAuth,
	}

	tcb := &bot.Bot{
		TwitchIRC: twitchIRC,
		Mongo:     mongoConnection,
		Helix:     helixClient,
		Logins:    make(map[string]string),
		Channels:  initChannels(ctx, mongoConnection, twitchIRC),
		Commands:  make(map[string]*bot.Command),
		Self:      self,
		StartTime: time.Now(),
	}

	// init actions that require bot.Bot object initialized already
	initializeEvents(tcb)

	err = tcb.TwitchIRC.Connect()
	if err != nil {
		log.Fatalln(err)
	}
}
