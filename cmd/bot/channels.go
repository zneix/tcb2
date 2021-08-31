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
	"github.com/zneix/tcb2/pkg/utils"
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

var channelSubscriptions = []eventsub.ChannelSubscription{
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
	// Fetch channel information for all channels (in bulks of 100)
	channelIDs := make([]string, 0, len(tcb.Channels))
	for k := range tcb.Channels {
		channelIDs = append(channelIDs, k)
	}

	channelIDChunks := utils.ChunkStringSlice(channelIDs, 100)

	for _, chunk := range channelIDChunks {
		go handleChannelsChunk(tcb, chunk)
	}
}

// handleChannelsChunk performs startup actions for channels with IDs from chunk
func handleChannelsChunk(tcb *bot.Bot, chunk []string) {
	// Fetch Title & Game for all loaded channels
	respC, err := tcb.Helix.GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: chunk,
	})
	if err != nil {
		log.Printf("Failed to query channel chunk %s; channels: %v", err, chunk)
		return
	}

	for _, respChannel := range respC.Data.Channels {
		channel := tcb.Channels[respChannel.BroadcasterID]

		// Set the ID in map translating login names back to IDs
		tcb.Logins[channel.Login] = channel.ID

		// Assign fetched properties to the channels
		channel.CurrentGame = respChannel.GameName
		channel.CurrentTitle = respChannel.Title

		// JOIN the channel
		tcb.TwitchIRC.Join(channel.Login)

		// Create all EventSub subscriptions parallelly
		for _, subscription := range channelSubscriptions {
			go func(sub eventsub.ChannelSubscription) {
				sub.ChannelID = channel.ID
				err = tcb.EventSub.CreateChannelSubscription(tcb.Helix, &sub)
				if err != nil {
					log.Println("[EventSub] Failed to create a subscription: " + err.Error())
				}
			}(subscription)
		}
	}

	// Fetch live status to check which channels out of loaded ones are live
	respS, err := tcb.Helix.GetStreams(&helix.StreamsParams{
		First:   100,
		UserIDs: chunk,
	})
	if err != nil {
		log.Printf("Failed to query stream chunk %s; channels: %v", err, chunk)
		return
	}

	for i := range respS.Data.Streams {
		respStream := &respS.Data.Streams[i]
		channel := tcb.Channels[respStream.UserID]

		log.Printf("[Helix:GetStreams] %s %s\n", respStream.Type, channel)
		if respStream.Type == "live" {
			channel.IsLive = true
		}
	}
}
