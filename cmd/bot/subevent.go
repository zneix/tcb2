package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// subEventTrigger will fetch relevant subscriptions and prepare ping messages, then attempt sending them in the channel where the event has occured
func subEventTrigger(msg *bot.SubEventMessage) {
	channel := msg.Bot.Channels[msg.ChannelID]
	// TODO: Change this once bot.Bot.Channels becomes a proper ChannelController
	if channel == nil {
		channel = msg.Bot.Self.Channel
	}

	cur, err := msg.Bot.Mongo.CollectionSubs(msg.ChannelID).Find(context.TODO(), bson.M{
		"event": msg.Type,
	})
	if err != nil {
		log.Printf("[Mongo] Failed querying events: " + err.Error())
		return
	}

	subs := []*bot.SubEventSubscription{}

	// value is either new title or new game depending of the event
	var value string
	switch msg.Type {
	case bot.SubEventTypeGame:
		value = channel.CurrentGame
	case bot.SubEventTypeTitle:
		value = channel.CurrentTitle
	}

	// Fetch all relevant subscriptions
	for cur.Next(context.TODO()) {
		// Deserialize sub data
		var sub *bot.SubEventSubscription
		err := cur.Decode(&sub)
		if err != nil {
			log.Println("[Mongo] Malformed subscription document: " + err.Error())
			continue
		}

		// Ignore subscriptions based on the value (if its present)
		if msg.Type == bot.SubEventTypeGame || msg.Type == bot.SubEventTypeTitle {
			if sub.Value != "" && !strings.Contains(value, sub.Value) {
				continue
			}
		}
		subs = append(subs, sub)
	}

	if len(subs) == 0 {
		log.Printf("[SubEvent] No relevant subscriptions for %v in %s\n", msg.Type, channel)
		return
	}

	// Construct ping message's "prefix"
	// Limit the length of a title / game in case it's too long, Twitch's limit is 140 anyway
	messagePrefix := strings.ReplaceAll(channel.Events[msg.Type], "{value}", utils.LimitString(value, 100))

	// Prepend the "[#channel]" part if the message is redirected from a channel
	// that is currently live and has EventsOnlyOffline flag set to true
	if channel.IsLive && channel.EventsOnlyOffline && msg.Type != bot.SubEventTypeLive {
		messagePrefix = fmt.Sprintf("[#%s] ", channel.Login)
	}
	messagePrefix = ".me " + messagePrefix

	// Prepare ping messages
	msgsToSend := []string{messagePrefix}
	for i := 0; i < len(subs); i++ {
		sub := subs[i]
		newMsg := fmt.Sprintf("%s %s", msgsToSend[len(msgsToSend)-1], sub.UserLogin)

		// Adding user to the message would exceed limit in the taget channel
		// We also want to re-run this event by decreasing i
		if utf8.RuneCountInString(newMsg) > channel.MessageLengthMax() {
			// We can't append any username to a message that is just our messagePrefix
			// Loop has to be broken or otherwise it'll run forever
			if msgsToSend[len(msgsToSend)-1] == messagePrefix {
				log.Println(fmt.Sprintf("[SubEvent] messagePrefix might be too long (%d) in %s: %# v", utf8.RuneCountInString(messagePrefix), channel, messagePrefix))
				break
			}
			msgsToSend = append(msgsToSend, messagePrefix)
			i--
			continue
		}
		// Otherwise it's good to append the username to message with pings
		msgsToSend[len(msgsToSend)-1] = newMsg
	}

	// Send messages to the target channel
	// TODO: Pajbot API (?)
	for i, v := range msgsToSend {
		log.Printf("[SubEvent] Announcing %s in %s; %d/%d(%d/%d chars)\n", msg.Type, channel, i+1, len(msgsToSend), utf8.RuneCountInString(v), channel.MessageLengthMax())
		channel.Send(v)
	}
}
