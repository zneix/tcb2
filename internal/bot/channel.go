package bot

import (
	"fmt"
	"log"
	"time"
)

func (channel *Channel) startMessageQueue(bot *Bot) {
	log.Println("Starting message queue for " + channel.String())
	defer log.Println("Done with message queue for " + channel.String())

	for message := range channel.QueueChannel {
		// Actually send the message to the chat
		bot.TwitchIRC.Say(channel.Login, message.Message)

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
