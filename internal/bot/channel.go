package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/zneix/tcb2/internal/mongo"
	"github.com/zneix/tcb2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (channel *Channel) StartMessageQueue(twitchIRC *twitch.Client) {
	log.Println("Starting message queue for " + channel.String())
	defer log.Println("[Channel] Message queue suddenly quit(?) for " + channel.String())

	for message := range channel.QueueChannel {
		// Actually send the message to the chat
		twitchIRC.Say(channel.Login, message.Message)

		// Update last sent message
		channel.LastMsg = message.Message

		// Wait for the cooldown
		time.Sleep(channel.Mode.MessageRatelimit())
	}
}

// Sends a message to the channel while taking care of ratelimits
func (channel *Channel) Send(message string) {
	// Don't attempt to send an empty message
	if message == "" {
		return
	}

	// TODO: Restrict usage of some commands, e.g. .ban, .timeout, .clear

	// limitting message length to not get it dropped
	message = utils.LimitString(message, channel.MessageLengthMax())

	// Append magic character at the end of the message if it's a duplicate
	if channel.LastMsg == message {
		message += " \U000E0000"
	}

	// Send message object to the message queue sending messages in ratelimit
	channel.QueueChannel <- &QueueMessage{
		Message: message,
	}
}

// Sendf formats according to a format specifier and runs channel.Send with the resulting string
func (channel *Channel) Sendf(format string, a ...interface{}) {
	channel.Send(fmt.Sprintf(format, a...))
}

func (channel *Channel) String() string {
	return fmt.Sprintf("#%s(%s)", channel.Login, channel.ID)
}

func (channel *Channel) MessageLengthMax() int {
	if channel.MessageLengthLimit > 0 {
		return channel.MessageLengthLimit
	}

	if channel.Mode == ChannelModeModerator {
		// Leaving 2 characters for the magic character
		return 498
	}
	// TODO: Investigate the actual limit for "pleb" modes (?)
	// mm2pl: maybe it's something like max of count(CHAR) / len(msg) for each unique character used in a message
	// mm2pl: or maybe it's some kind of GOW average
	// mm2pl: max((msg.count(ch) / len(msg) for ch in set(msg))) seems like a good approximation
	// For now I'm lazy and just gonna hardcode some reasonable value in here
	return 468
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
