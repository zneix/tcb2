package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
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
	log.Println("Starting titlechange_bot", common.Version())

	cfg := config.New()
	ctx := context.Background()

	mongoConnection := mongo.NewMongoConnection(ctx, cfg)
	mongoConnection.Connect(ctx)

	// twitch read conn
	twitchRead := twitch.NewAnonymousClient()
	twitchRead.SetJoinRateLimiter(twitch.CreateVerifiedRateLimiter())

	// twitch write conn
	twitchWrite := twitch.NewClient(cfg.TwitchLogin, "oauth:"+cfg.TwitchOAuth)

	helixClient, err := helixclient.New(cfg)
	if err != nil {
		log.Fatalln("[Helix] Error while initializing client:", err)
	}

	apiServer := api.New(cfg)

	esub := eventsub.New(cfg, apiServer)

	tcb := &bot.Bot{
		TwitchRead:  twitchRead,
		TwitchWrite: twitchWrite,
		Mongo:       mongoConnection,
		Helix:       helixClient,
		EventSub:    esub,
		Logins:      make(map[string]string),
		Channels:    loadChannels(ctx, mongoConnection, twitchWrite),
		Commands:    bot.NewCommandController(cfg.CommandPrefix),
		Self: &bot.Self{
			Login: cfg.TwitchLogin,
			OAuth: cfg.TwitchOAuth,
		},
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

	// Connect twitch connections
	// TODO: Use proper waiting for both connections - maybe use channels waiting for closure
	// Check twitch connection's .Connect for a good example
	// For now as a scuffed fix read connection will be blocking
	go func() {
		err = tcb.TwitchWrite.Connect()
		if err != nil {
			log.Fatalln(fmt.Errorf("twitch write connection errored: %w", err))
		}
	}()

	err = tcb.TwitchRead.Connect()
	if err != nil {
		log.Fatalln(fmt.Errorf("twitch read connection errored: %w", err))
	}
}
