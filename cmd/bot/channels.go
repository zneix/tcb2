package main

import (
	"context"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/eventsub"
	"github.com/zneix/tcb2/internal/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// loadChannels fetches configured channels from the database, sets default values and message queue for each of them
func loadChannels(bgctx context.Context, mongoConn *mongo.Connection, twitchIRC *twitch.Client) map[string]*bot.Channel {
	channels := make(map[string]*bot.Channel)

	ctx, cancel := context.WithTimeout(bgctx, 10*time.Second)
	defer cancel()

	// Query all channels from the database, excluding inactive channels
	cur, err := mongoConn.Collection(mongo.CollectionNameChannels).Find(ctx, bson.M{
		"mode": &bson.M{
			"$ne": int(bot.ChannelModeInactive),
		},
	})
	if err != nil {
		log.Fatalln("[Mongo] Error querying channels: " + err.Error())
	}

	for cur.Next(ctx) {
		// Deserialize channel data
		var channel bot.Channel
		err := cur.Decode(&channel)
		if err != nil {
			log.Println("[Mongo] Malformed channel document: " + err.Error())
			continue
		}

		// Initialize default values
		channel.QueueChannel = make(chan *bot.QueueMessage)
		go channel.StartMessageQueue(twitchIRC)

		channels[channel.ID] = &channel
	}

	if err := cur.Err(); err != nil {
		log.Println("[Mongo] Last cursor error while loading channels wasn't nil: " + err.Error())
	}

	return channels
}

// joinChannels performs startup actions for all the channels that are already loaded
func joinChannels(tcb *bot.Bot) {
	for ID, channel := range tcb.Channels {
		// Set the ID in map translating login names back to IDs
		tcb.Logins[channel.Login] = ID

		// JOIN the channel
		tcb.TwitchIRC.Join(channel.Login)

		// Create EventSub subscriptions (but don't block)
		go func(channelID string) {
			err := tcb.EventSub.CreateChannelSubscription(tcb.Helix, &eventsub.ChannelSubscription{
				Type:      "channel.update",
				Version:   "1",
				ChannelID: channelID,
			})
			if err != nil {
				log.Println("[EventSub] Failed to create a subscription: " + err.Error())
			}
			err = tcb.EventSub.CreateChannelSubscription(tcb.Helix, &eventsub.ChannelSubscription{
				Type:      "stream.online",
				Version:   "1",
				ChannelID: channelID,
			})
			if err != nil {
				log.Println("[EventSub] Failed to create a subscription: " + err.Error())
			}
			err = tcb.EventSub.CreateChannelSubscription(tcb.Helix, &eventsub.ChannelSubscription{
				Type:      "stream.offline",
				Version:   "1",
				ChannelID: channelID,
			})
			if err != nil {
				log.Println("[EventSub] Failed to create a subscription: " + err.Error())
			}
		}(channel.ID)
	}
}
