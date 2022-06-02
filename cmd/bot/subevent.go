package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/mongo"
	"github.com/zneix/tcb2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"
)

// subEventTrigger will fetch relevant subscriptions and prepare ping messages, then attempt sending them in the channel where the event has occured
func subEventTrigger(msg *bot.SubEventMessage) {
	channel := msg.Bot.Channels[msg.ChannelID]
	ctx := context.TODO()

	curSubs, err := msg.Bot.Mongo.CollectionSubs(msg.ChannelID).Find(ctx, bson.M{
		"event": msg.Type,
	})
	if err != nil {
		log.Println("[Mongo] Failed querying events:", err)
		return
	}

	subs := make([]*bot.SubEventSubscription, 0) // XXX: Test if this won't panic

	// value is either new title or new game depending of the event
	var value string
	switch msg.Type {
	case bot.SubEventTypeGame:
		value = channel.CurrentGame
	case bot.SubEventTypeTitle:
		value = channel.CurrentTitle
	}

	valueLower := strings.ToLower(value)

	// Fetch all relevant subscriptions
	for curSubs.Next(ctx) {
		// Deserialized sub data
		sub := new(bot.SubEventSubscription)
		err := curSubs.Decode(&sub)
		if err != nil {
			log.Println("[Mongo] Malformed subscription document:", err)
			continue
		}

		// Ignore subscriptions based on the value (if its present)
		if msg.Type == bot.SubEventTypeGame || msg.Type == bot.SubEventTypeTitle {
			if sub.Value != "" && !strings.Contains(valueLower, strings.ToLower(sub.Value)) {
				continue
			}
		}
		subs = append(subs, sub)
	}

	if len(subs) == 0 {
		log.Printf("[SubEvent] No relevant subscriptions for %v in %s\n", msg.Type, channel)
		return
	}

	// If the message should be redirected from a channel that is currently live
	// and has EventsOnlyOffline flag set to true,
	redirect := false

	if channel.IsLive && channel.EventsOnlyOffline && msg.Type != bot.SubEventTypeLive {
		log.Printf("[SubEvent] Redirecting %s from %s because channel is live\n", msg.Type, channel)
		redirect = true
	}

	messagePrefix := createMessagePrefix(channel.Events[msg.Type], value, channel.Login, redirect)

	// Prepare ping messages
	msgsToSend := []string{messagePrefix}
	for i := 0; i < len(subs); {
		sub := subs[i]
		newMsg := fmt.Sprintf("%s %s", msgsToSend[len(msgsToSend)-1], sub.UserLogin)

		// Adding user to the message would exceed limit in the target channel
		// We also want to re-run this iteration by decreasing i
		if utf8.RuneCountInString(newMsg) > channel.MessageLengthMax() {
			// We can't append any username to a message that is just our messagePrefix
			// Loop has to be broken or otherwise it'll run forever
			if msgsToSend[len(msgsToSend)-1] == messagePrefix {
				log.Printf("[SubEvent] messagePrefix might be too long (%d) in %s: %# v\n", utf8.RuneCountInString(messagePrefix), channel, messagePrefix)
				break
			}
			msgsToSend = append(msgsToSend, messagePrefix)
			continue
		}
		// Otherwise it's good to append the username to message with pings
		msgsToSend[len(msgsToSend)-1] = newMsg
		i++
	}

	// In case of EventsOnlyOffline flag, send messages to bot's own channel
	if redirect {
		botID, ok := msg.Bot.Logins[msg.Bot.Self.Login]
		if !ok {
			// handle error
			return
		}

		channel, ok = msg.Bot.Channels[botID]
		if !ok {
			// handle error as well
			return
		}
	}

	// TODO: Pajbot API (?)
	for i, v := range msgsToSend {
		log.Printf("[SubEvent] Announcing %s in %s; %d/%d(%d/%d chars)\n", msg.Type, channel, i+1, len(msgsToSend), utf8.RuneCountInString(v), channel.MessageLengthMax())
		channel.Send(v)
	}

	// Fetch and send channel's MOTD
	handleMOTD(msg)
}

// createMessagePrefix constructs ping message's "prefix", which will be the beginning of every ping message
func createMessagePrefix(format, value, login string, redirect bool) string {
	// Limit the length of a title / game in case it's too long, Twitch's limit is 140 anyway
	prefixReplacer := strings.NewReplacer(
		"{value}", utils.LimitString(value, 100),
		"{login}", login,
	)

	prefix := ".me "

	if redirect {
		prefix += fmt.Sprintf("[#%s] ", login)
	}

	return prefix + prefixReplacer.Replace(format)
}

// handleMOTD queries MOTD for the channel where event occurred and sends it if exists
func handleMOTD(msg *bot.SubEventMessage) {
	// By design, it is only sent on live event
	if msg.Type != bot.SubEventTypeLive {
		return
	}

	res := msg.Bot.Mongo.Collection(mongo.CollectionNameMOTDs).FindOne(context.TODO(), bson.M{
		"channel_id": msg.ChannelID,
	})
	if err := res.Err(); err != nil {
		// Ignoring ErrNoDocuments, since it's not really an error (plus it is returned quite often)
		if !errors.Is(err, mongodb.ErrNoDocuments) {
			log.Printf("[Mongo] Failed querying MOTD for %s: %s\n", msg.ChannelID, err)
		}
		return
	}

	// Deserialized MOTD data
	motd := new(bot.SubEventMOTD)
	err := res.Decode(&motd)
	if err != nil {
		log.Printf("[Mongo] Malformed MOTD document for %s: %s\n", msg.ChannelID, err)
		return
	}

	// Send the MOTD to the channel
	channel := msg.Bot.Channels[msg.ChannelID]
	log.Printf("[SubEvent] Sending channel MOTD to %s (%d/%d chars)\n", channel, utf8.RuneCountInString(motd.Message), channel.MessageLengthMax())
	// Channel.Send handles empty messages for us already, in case motd.Message would be empty
	channel.Send(motd.Message)
}
