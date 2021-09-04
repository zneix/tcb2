package commands

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NotifyMe(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "notifyme",
		Aliases:         []string{"tcbnotifyme"},
		Description:     "Subscribe to an event. Optional value can be only used with title and game events. For list of available events use: {prefix}events",
		Usage:           "<event> [optional value]",
		CooldownChannel: 1 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]

			eventStrings := []string{}
			for i, desc := range bot.SubEventDescriptions {
				eventStrings = append(eventStrings, fmt.Sprintf("%s (%s)", bot.SubEventType(i), desc))
			}
			availableEvents := strings.Join(eventStrings, ", ")

			// No arguments, return an error message
			if len(args) < 1 {
				channel.Sendf("@%s, you must specify an event to subscribe to. Available events: %s", msg.User.Name, availableEvents)
				return
			}

			// Parse sub event type passed as the first argument
			valid, event := bot.ParseSubEventType(strings.ToLower(args[0]))
			if !valid {
				channel.Sendf("@%s, given event name is not valid. %s", msg.User.Name, "TODO: show all events you can subscribe to")
				return
			}

			// Determine the optional value
			value := strings.Join(args[1:], " ")
			if event != bot.SubEventTypeGame && event != bot.SubEventTypeTitle {
				value = ""
			}

			// Find user's subscriptions in this chat for the specified event
			cur, err := tcb.Mongo.CollectionSubs(msg.RoomID).Find(context.TODO(), bson.M{
				"user_id": msg.User.ID,
				"event":   event,
			})
			if err != nil {
				log.Println("[Mongo] Failed querying subscription: " + err.Error())
				return
			}

			hasThisSub := false
			hasThisSubAllValues := false
			hasThisSubWithThisValue := false
			var deletedSubCount int

			// Deserialize subscription data
			for cur.Next(context.TODO()) {
				var sub *bot.SubEventSubscription
				err = cur.Decode(&sub)
				if err != nil {
					log.Println("[Mongo] Malformed subscription document: " + err.Error())
					continue
				}

				if sub.Event != event {
					continue
				}

				hasThisSub = true
				if sub.Value == "" {
					hasThisSubAllValues = true
				}
				if strings.EqualFold(sub.Value, value) {
					hasThisSubWithThisValue = true
				}
			}

			// User already has a subscription with the exact value
			// inform them about it and how can they unsubscribe
			if hasThisSubWithThisValue {
				reply := fmt.Sprintf("@%s, you already have a subscription to the event %s", msg.User.Name, event)
				if len(value) > 0 {
					reply += " with the provided value FeelsDankMan .."
				}

				channel.Sendf("%s. If you want to unsubscribe, use: %sremoveme %s", reply, tcb.Commands.Prefix, event)
				return
			}

			// User has a subscription for this event for all values
			if hasThisSubAllValues && len(value) > 0 {
				channel.Sendf("@%s, you already have a subscription for event %s that matches all values. If you want to be pinged only on specific values, use \"%sremoveme %s\" first before running this command again", msg.User.Name, event, tcb.Commands.Prefix, event)
				return
			}

			// User has subscription(s) for this event, but requests a subscription for all values
			// We want to delete all other subscriptions for this event for this user first
			if hasThisSub && len(value) < 1 {
				var res *mongo.DeleteResult
				res, err = tcb.Mongo.CollectionSubs(msg.RoomID).DeleteMany(context.TODO(), bson.M{
					"user_id": msg.User.ID,
					"event":   event,
				})
				if err != nil {
					log.Println("[Mongo] Failed deleting subscriptions: " + err.Error())
					channel.Sendf("@%s, internal server error occured while trying to delete your old subscriptions monkaS @zneix", msg.User.Name)
					return
				}
				deletedSubCount = int(res.DeletedCount)
				log.Printf("[Mongo] Deleted %d subscription(s) for %# v(%s) in %s", res.DeletedCount, msg.User.Name, msg.User.ID, channel)
			}

			// Add requested subscription
			res, err := tcb.Mongo.CollectionSubs(msg.RoomID).InsertOne(context.TODO(), bot.SubEventSubscription{
				UserID:    msg.User.ID,
				UserLogin: msg.User.Name,
				Event:     event,
				Value:     value,
			})
			if err != nil {
				log.Println("[Mongo] Failed adding new subscription: " + err.Error())
				channel.Sendf("@%s, internal server error occured while trying to add your new subscription monkaS @zneix", msg.User.Name)
				return
			}
			log.Printf("[Mongo] Added 1 subscription for %# v(%s) in %s, ID: %v", msg.User.Name, msg.User.ID, channel, res.InsertedID)

			var givenValue string
			if len(value) > 0 {
				givenValue = fmt.Sprintf(", but only when the %s contains provided value", event)
			}
			reply := fmt.Sprintf("@%s, I will now ping you when the %s%s!", msg.User.Name, bot.SubEventDescriptions[event], givenValue)
			if hasThisSub && len(value) < 1 {
				// We had to remove all other subscriptions for this event and add a new one
				reply += fmt.Sprintf(" You previously had %d subscription(s) for this event that were set to only match specific values. These subscriptions have been removed and you will now be notified regardless of the value SeemsGood", deletedSubCount)
			}

			// Return the final message to the user
			channel.Send(reply)
		},
	}
}
