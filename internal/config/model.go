package config

type TCBConfig struct {
	/// Bot

	// CommandPrefix prefix to which bot will respond to (chat commands)
	CommandPrefix string `yaml:"command-prefix"`

	/// Misc

	// SupinicAPIKey bot's key to the Supinic's API, in the format: "SupibotID:APIKey" without quotation marks
	SupinicAPIKey string `yaml:"supinic-api-key"`

	/// API

	// BaseURL url of the API to which clients will make their requests. Useful if the API is proxied through reverse proxy like nginx.
	// Its value needs to contain full URL with protocol scheme, e.g. https://braize.pajlada.com/chatterino
	BaseURL string `yaml:"base-url"`
	// BindAddress address to which API will bind and start listening on
	BindAddress string `yaml:"bind-address"`

	/// Twitch

	// TwitchLogin login name of the account on which bot will log into Twitch's IRC on the Write connection
	TwitchLogin string `yaml:"twitch-login"`
	// TwitchOAuth OAuth token of the account on which bot will log into Twitch's IRC on the Write connection
	// It should not have the "oauth:" prefix - that is added once bot attempts to authenticate
	TwitchOAuth string `yaml:"twitch-oauth"`
	// TwitchClientID
	TwitchClientID string `yaml:"twitch-client-id"`
	// TwitchClientSecret
	TwitchClientSecret string `yaml:"twitch-client-secret"`
	// TwitchEventSubSecret secret used to create subscriptions and verify incoming notifications
	// Must be between 10 and 100 characters long
	TwitchEventSubSecret string `yaml:"twitch-eventsub-secret"`

	/// Mongo 🥭

	// MongoUsername
	MongoUsername string `yaml:"mongo-username"`
	// MongoPassword
	MongoPassword string `yaml:"mongo-password"`
	// MongoPort port to which connection will try to connect
	// Host is hardcoded to localhost due to security concerns (use ssh port tunneling while testing/developing on a remote machine)
	MongoPort string `yaml:"mongo-port"`
	// MongoDatabaseNamename name of the database used by the bot
	// It's treated as a prefix to allow using separate databases while testing/developing
	MongoDatabaseName string `yaml:"mongo-database-name"`
	// MongoAuthDB name of authentication databse, used as AuthSource while creating a new mongo.Connection
	// This should usually be left unchanged
	MongoAuthDB string `yaml:"mongo-auth-db"`
}
