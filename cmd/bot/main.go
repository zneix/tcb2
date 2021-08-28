package main

import (
	"context"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/api"
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/config"
	"github.com/zneix/tcb2/internal/eventsub"
	"github.com/zneix/tcb2/internal/helixclient"
	"github.com/zneix/tcb2/internal/mongo"
)

const (
	VERSION        = "2.0-alpha"
	COMMAND_PREFIX = "!"
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

	apiServer := api.New(cfg)

	esub := eventsub.New(cfg, apiServer)

	self := &bot.Self{
		Login: cfg.TwitchLogin,
		OAuth: cfg.TwitchOAuth,
	}

	tcb := &bot.Bot{
		TwitchIRC: twitchIRC,
		Mongo:     mongoConnection,
		Helix:     helixClient,
		EventSub:  esub,
		Logins:    make(map[string]string),
		Channels:  initChannels(ctx, mongoConnection, twitchIRC),
		Commands:  bot.NewCommandController(),
		Self:      self,
		StartTime: time.Now(),
	}

	// Init actions that require bot.Bot object initialized already
	initializeEvents(tcb)
	registerCommands(tcb)

	// TODO: Manage goroutines below and (currently blocking) Connect() with sync.WaitGroup
	// Listen on the API instance
	go apiServer.Listen()

	err = tcb.TwitchIRC.Connect()
	if err != nil {
		log.Fatalln(err)
	}
}
