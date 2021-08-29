package main

import (
	"context"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/nicklaw5/helix"
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

var channelSubscriptions = []*eventsub.ChannelSubscription{
	{
		Type:    "channel.update",
		Version: "1",
	},
	{
		Type:    "stream.online",
		Version: "1",
	},
	{
		Type:    "stream.offline",
		Version: "1",
	},
}

// joinChannels performs startup actions for all the channels that are already loaded
func joinChannels(tcb *bot.Bot) {
	// TODO: Fetch channel information for all channels (in bulks of 100)
	//var channelChunks [][]string
	for ID := range tcb.Channels {
		tcb.Helix.GetChannelInformation(&helix.GetChannelInformationParams{
			BroadcasterID: ID,
		})
	}

	for ID, channel := range tcb.Channels {
		// Set the ID in map translating login names back to IDs
		tcb.Logins[channel.Login] = ID

		// JOIN the channel
		tcb.TwitchIRC.Join(channel.Login)

		// TODO: Assign data from stuff fetched earlier to the channel

		// Putting API-based actions to a goroutine to make parallel loading faster
		go func(channelID string) {
			// Create all EventSub subscriptions parallelly
			for _, subscription := range channelSubscriptions {
				go func(sub *eventsub.ChannelSubscription) {
					sub.ChannelID = channelID
					err := tcb.EventSub.CreateChannelSubscription(tcb.Helix, sub)
					if err != nil {
						log.Println("[EventSub] Failed to create a subscription: " + err.Error())
					}
				}(subscription)
			}
		}(channel.ID)
	}
}
