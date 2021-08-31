package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/zneix/tcb2/internal/bot"
	"go.mongodb.org/mongo-driver/bson"
)

// subEventTrigger will fetch relevant subscriptions and prepare ping messages, then attempt sending them in the channel where the event has occured
func subEventTrigger(msg *bot.SubEventMessage) {
	cur, err := msg.Bot.Mongo.CollectionSubs(msg.ChannelID).Find(context.TODO(), bson.M{
		"event": msg.Type,
	})
	if err != nil {
		log.Printf("[Mongo] Failed querying events: " + err.Error())
		return
	}

	subs := []*bot.SubEventSubscription{}
	channel := msg.Bot.Channels[msg.ChannelID]

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
	messagePrefix := ".me " + strings.ReplaceAll(channel.Events[msg.Type], "{value}", value)

	// Prepare ping messages
	msgsToSend := []string{messagePrefix}
	for _, sub := range subs {
		lastMsg := msgsToSend[len(msgsToSend)-1]
		newMsg := fmt.Sprintf("%s %s", lastMsg, sub.UserLogin)

		// Adding user to the message would exceed limit in the taget channel
		if len(newMsg) > channel.MessageLengthMax() {
			msgsToSend = append(msgsToSend, messagePrefix)
			continue
		}
		// Otherwise it's good to append the username to message with pings
		msgsToSend[len(msgsToSend)-1] = newMsg
	}

	// Send messages to the target channel
	// TODO: Pajbot API (?)
	for i, v := range msgsToSend {
		log.Printf("[SubEvent] Announcing %s in %s, %d/%d(%d chars)\n", msg.Type, channel, i+1, len(msgsToSend), len(v))
		channel.Send(v)
	}
}
