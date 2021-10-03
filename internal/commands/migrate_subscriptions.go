package commands

import (
	"context"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/zneix/tcb2/internal/bot"
	"github.com/zneix/tcb2/internal/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type SubEventSubscriptionOld struct {
	Channel string `bson:"channel"`
	User    string `bson:"user"`
	Event   string `bson:"event"`
	Value   string `bson:"requiredValue"`
}

func MigrateSubscriptions(tcb *bot.Bot) *bot.Command {
	return &bot.Command{
		Name:            "migrate_subscriptions",
		Aliases:         []string{},
		Description:     "Migration command, admin use only",
		Usage:           "",
		IgnoreSelf:      false,
		CooldownChannel: 3 * time.Second,
		CooldownUser:    5 * time.Second,
		Run: func(msg twitch.PrivateMessage, args []string) {
			if msg.User.Name != "zneix" {
				return
			}

			channel := tcb.Channels[msg.RoomID]

			// Helix username cache: User Login -> User ID
			xd := make(map[string]string)
			errCount := 0

			// Iterate over all channels and perform migrations for these
			channel.Sendf("Attempting to migrate subscriptions for %d channels", len(tcb.Channels))
			for _, c := range tcb.Channels {
				channel.Sendf("Migrating subscriptions for %s", c)

				cur, err := tcb.Mongo.Collection(mongo.CollectionNameOldSubscriptions).Find(context.Background(), bson.M{
					"channel": c.Login,
				})
				if err != nil {
					errCount++
					log.Printf("[Mongo] Error querying old subscriptions for %s: %s\n", c, err)
					channel.Sendf("Failed to migrate subscriptions for %s: %s", c, err)
					continue
				}

				for cur.Next(context.TODO()) {
					// Deserialize old subscription data
					var subOld *SubEventSubscriptionOld
					err = cur.Decode(&subOld)
					if err != nil {
						errCount++
						log.Println("[Mongo] Malformed old subscription document: " + err.Error())
						continue
					}

					userID, ok := xd[subOld.User]
					if !ok {
						// When the userID wasn't found in cache, query helix for it
						log.Printf("[Helix] Querying %s\n", subOld.User)
						res, err := tcb.Helix.GetUsers(&helix.UsersParams{
							Logins: []string{subOld.User},
						})
						if err != nil {
							errCount++
							log.Printf("[Helix] Failed to query user %s: %s\n", subOld.User, err)
							channel.Sendf("Failed to query user %s: %s", subOld.User, err)
							continue
						}

						if len(res.Data.Users) != 1 {
							log.Printf("User %s is banned\n", subOld.User)
							userID = "banned"
						} else {
							userID = res.Data.Users[0].ID
						}
						xd[subOld.User] = userID
					}

					valid, event := bot.ParseSubEventType(subOld.Event)
					if !valid {
						errCount++
						log.Printf("Invalid event on old subscription: %# v\n", subOld)
						continue
					}
					subNew := &bot.SubEventSubscription{
						UserLogin: subOld.User,
						UserID:    userID,
						Event:     event,
						Value:     subOld.Value,
					}
					resIns, err := tcb.Mongo.CollectionSubs(c.ID).InsertOne(context.TODO(), *subNew)
					if err != nil {
						errCount++
						log.Printf("[Mongo] Error inserting new subscription %s; %# v\n", err, subNew)
						continue
					}
					log.Printf("[Mongo] Inserted new subscription %# v; ID: %v\n", subNew, resIns.InsertedID)
				}
			}
			channel.Sendf("Finished migrating subscriptions for all %d channels KKona encountered %d errors @zneix", len(tcb.Channels), errCount)
		},
	}
}
