package commands

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/bot"
)

func RemoveMe(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "removeme",
		Aliases:         []string{"tcbremoveme"},
		Description:     "Subscribe to an event. Optional value can be only used with title and game events. For list of available events use: {prefix}events",
		Usage:           "<event> [optional value]",
		CooldownChannel: 1 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]

			// Parts of responses used commonly across the command
			checkAllEvents := fmt.Sprintf("Check all events you've subscribed to with: %ssubscribed", tcb.Commands.Prefix)

			// No arguments, return an error message
			if len(args) < 1 {
				channel.Sendf("@%s, you must specify an event to unsubscribe from. %s", msg.User.Name, checkAllEvents)
				return
			}

			// Parse sub event type passed as the first argument
			valid, event := bot.ParseSubEventType(strings.ToLower(args[0]))
			if !valid {
				channel.Sendf("@%s, given event name is not valid. %s", msg.User.Name, checkAllEvents)
				return
			}

			// Determine the optional value
			value := strings.Join(args[1:], " ")
			if event != bot.SubEventTypeGame && event != bot.SubEventTypeTitle {
				value = ""
			}

			// If value is empty, remove the user from all subscriptions to this event
			removeQuery := &bot.SubEventSubscription{
				UserID: msg.User.ID,
				Event:  event,
			}
			// Otherwise, only remove subscriptions that match that value (case sensitive)
			if value != "" {
				removeQuery.Value = value
			}
			res, err := tcb.Mongo.CollectionSubs(msg.RoomID).DeleteMany(context.TODO(), removeQuery)
			if err != nil {
				log.Println("[Mongo] Failed deleting subscriptions: " + err.Error())
				channel.Sendf("@%s, internal server error occured while trying to delete your subscriptions monkaS @zneix", msg.User.Name)
				return
			}
			fmt.Printf("%# v\n", res)
			log.Printf("[Mongo] Deleted %d subscription(s) for %# v(%s) in %s", res.DeletedCount, msg.User.Name, msg.User.ID, channel)

			if res.DeletedCount == 0 {
				if len(value) > 0 {
					// Didn't match the value
					channel.Sendf("@%s, you are not subscribed to the event %s with provided value FeelsDankMan %s", msg.User.Name, event, checkAllEvents)
				} else {
					// User wasn't subscribed to this event
					channel.Sendf("@%s, you are not subscribed to the event %s. %s", msg.User.Name, event, checkAllEvents)
				}
				return
			}

			channel.Sendf("@%s, successfully removed %d subscription(s) to event %s", msg.User.Name, res.DeletedCount, event)
		},
	}
}
