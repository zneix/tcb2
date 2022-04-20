package main

import (
	"context"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/zneix/tcb2/internal/api"
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/common"
	"github.com/zneix/tcb2/internal/config"
	"github.com/zneix/tcb2/internal/eventsub"
	"github.com/zneix/tcb2/internal/helixclient"
	"github.com/zneix/tcb2/internal/mongo"
	"github.com/zneix/tcb2/internal/supinicapi"
)

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
}

func main() {
	log.Printf("Starting titlechange_bot %s", common.Version())

	cfg := config.New()
	ctx := context.Background()

	mongoConnection := mongo.NewMongoConnection(ctx, cfg)
	mongoConnection.Connect()

	twitchIRC := twitch.NewClient(cfg.TwitchLogin, "oauth:"+cfg.TwitchOAuth)
	twitchIRC.SetJoinRateLimiter(twitch.CreateVerifiedRateLimiter())

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
		Channels:  loadChannels(ctx, mongoConnection, twitchIRC),
		Commands:  bot.NewCommandController(cfg.CommandPrefix),
		Self:      self,
		StartTime: time.Now(),
	}

	// Register actions that require bot.Bot object initialized already
	registerEvents(tcb)
	registerCommands(tcb)

	// TODO: Manage goroutines below and (currently blocking) Connect() with sync.WaitGroup
	// Listen on the API instance
	go apiServer.Listen()

	// Ping Supinic's API periodically to signal that bot is alive
	supinic := supinicapi.New(cfg.SupinicAPIKey)
	go supinic.UpdateAliveStatus()

	err = tcb.TwitchIRC.Connect()
	if err != nil {
		log.Fatalln(err)
	}
}
