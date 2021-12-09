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
)

func Subscribed(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "subscribed",
		Aliases:         []string{"tcbsubscribed"},
		Description:     "Shows you list of events you're subscribed to",
		Usage:           "",
		CooldownChannel: 100 * time.Millisecond,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			channel := tcb.Channels[msg.RoomID]

			// Find all user's subscriptions in this chat
			// TODO: Consider adding a way to check your subscriptions in other / all chats
			cur, err := tcb.Mongo.CollectionSubs(msg.RoomID).Find(context.TODO(), bson.M{
				"user_id": msg.User.ID,
			})
			if err != nil {
				log.Printf("[Mongo] Failed querying events: %s\n", err)
				return
			}

			// subMap contains indexes of subscriptions in the subs slice
			subMap := make(map[bot.SubEventType][]int)
			subs := []*bot.SubEventSubscription{}

			// Fetch all relevant subscriptions
			for cur.Next(context.TODO()) {
				// Deserialize sub data
				var sub *bot.SubEventSubscription
				err := cur.Decode(&sub)
				if err != nil {
					log.Println("[Mongo] Malformed subscription document: " + err.Error())
					continue
				}
				subs = append(subs, sub)
				subMap[sub.Event] = append(subMap[sub.Event], len(subs)-1)
			}

			// User isn't subscribed to anything, tell them how can they do that
			if len(subs) == 0 {
				// @zneix, You are not subscribed to any events. Use !notifyme <event> [optional value] to subscribe. Valid events are: game, live, offline, title
				eventStrings := []string{}
				for i, desc := range bot.SubEventDescriptions {
					eventStrings = append(eventStrings, fmt.Sprintf("%s (%s)", bot.SubEventType(i), desc))
				}
				channel.Sendf("@%s, you are not subscribed to any events. Use %snotifyme to subscribe to an event. Valid events: %s", msg.User.Name, tcb.Commands.Prefix, strings.Join(eventStrings, ", "))
				return
			}

			// Inform the user about their subscriptions
			parts := []string{}
			for k, v := range subMap {
				values := []string{}
				for _, subIndex := range v {
					if subs[subIndex].Value != "" {
						values = append(values, fmt.Sprintf("%q", subs[subIndex].Value))
					}
				}
				if len(values) == 0 {
					parts = append(parts, k.String())
				} else {
					parts = append(parts, fmt.Sprintf("%s (only for values: %s)", k, strings.Join(values, ", ")))
				}
			}
			channel.Sendf("@%s, you have %d subscription(s) to %d event(s): %s", msg.User.Name, len(subs), len(subMap), strings.Join(parts, ", "))
		},
	}
}
