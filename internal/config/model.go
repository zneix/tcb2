package config

type TCBConfig struct {
	// API

	BaseURL     string `mapstructure:"base-url"`
	BindAddress string `mapstructure:"bind-address"`

	// Twitch

	TwitchLogin          string `mapstructure:"twitch-login"`
	TwitchOAuth          string `mapstructure:"twitch-oauth"`
	TwitchClientID       string `mapstructure:"twitch-client-id"`
	TwitchClientSecret   string `mapstructure:"twitch-client-secret"`
	TwitchEventSubSecret string `mapstructure:"twitch-eventsub-secret"`

	// Mongo 🥭

	MongoUsername     string `mapstructure:"mongo-username"`
	MongoPassword     string `mapstructure:"mongo-password"`
	MongoPort         string `mapstructure:"mongo-port"`
	MongoDatabaseName string `mapstructure:"mongo-database-name"`
	MongoAuthDB       string `mapstructure:"mongo-auth-db"`
}
