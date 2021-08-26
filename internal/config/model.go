package config

type TCBConfig struct {
	// API (Eventsub)

	BaseURL     string `mapstructure:"base-url"`
	BindAddress string `mapstructure:"bind-address"`

	// Twitch

	TwitchLogin        string `mapstructure:"twitch-login"`
	TwitchOAuth        string `mapstructure:"twitch-oauth"`
	TwitchClientID     string `mapstructure:"twitch-client-id"`
	TwitchClientSecret string `mapstructure:"twitch-client-secret"`

	// TODO: Mongo
}
