package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func (channel *Channel) StartMessageQueue(twitchIRC *twitch.Client) {
	log.Println("Starting message queue for " + channel.String())
	defer log.Println("Done with message queue for " + channel.String())

	for message := range channel.QueueChannel {
		// Actually send the message to the chat
		twitchIRC.Say(channel.Login, message.Message)

		// Update last sent message
		channel.LastMsg = message.Message

		// Wait for the cooldown
		time.Sleep(channel.Mode.MessageRatelimit())
	}
}

func (channel *Channel) Send(message string) {
	// Don't attempt to send an empty message
	if len(message) == 0 {
		return
	}

	// TODO: Restrict usage of some commands, e.g. .ban, .timeout, .clear

	// limitting message length to not get it dropped
	if len(message) > channel.Mode.MessageLengthMax() {
		message = message[0:channel.Mode.MessageLengthMax()-3] + "..."
	}

	// Append magic character at the end of the message if it's a duplicate
	if channel.LastMsg == message {
		message += " \U000E0000"
	}

	// Send message object to the message queue sending messages in ratelimit
	channel.QueueChannel <- &QueueMessage{
		Message: message,
	}
}

func (channel *Channel) String() string {
	return fmt.Sprintf("#%s(%s)", channel.Login, channel.ID)
}

func (channel *Channel) ChangeMode(mongoConn *mongo.Connection, newMode ChannelMode) (err error) {
	log.Printf("[Mongo] Changing mode in %s from %v to %v", channel.String(), channel.Mode, newMode)
	channel.Mode = newMode

	// Update mode in the database as well
	_, err = mongoConn.Collection(mongo.CollectionNameChannels).UpdateOne(context.TODO(), bson.M{
		"id": channel.ID,
	}, bson.M{
		"$set": bson.M{
			"mode": newMode,
		},
	})

	if err != nil {
		log.Printf("[Mongo] Error updating ChannelMode for %s: %s\n", channel.String(), err.Error())
	}
	return
}
