package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	appName    = "tcb2"
	configName = "config"
)

// readFromPath reads the config values from the given path (i.e. path/${configName}.yaml) and returns its values as a map.
// This allows us to use mergeConfig cleanly
func readFromPath(path string) (values map[string]interface{}, err error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	if err = v.ReadInConfig(); err != nil {
		notFoundError := &viper.ConfigFileNotFoundError{}
		if errors.As(err, notFoundError) {
			err = nil
			return
		}
		return
	}

	v.Unmarshal(&values)

	return
}

// mergeConfig uses viper.MergeConfigMap to read config values in the unix
// standard, so you start furthest down with reading the system config file,
// merge those values into the main config map, then read the home directory
// config files, and merge any set values from there, and lastly the config
// file in the cwd and merge those in. If a value is not set in the cwd config
// file, but one is set in the system config file then the system config file
// value will be used
func mergeConfig(v *viper.Viper, configPaths []string) {
	for _, configPath := range configPaths {
		configMap, err := readFromPath(configPath)
		if err != nil {
			fmt.Printf("Error reading config file from %s.yaml: %s\n", filepath.Join(configPath, configName), err)
			return
		}
		v.MergeConfigMap(configMap)
	}
}

func init() {
	// Define command-line flags and default values

	// API
	pflag.StringP("base-url", "b", "", "Base URL of the API to which clients will make their requests. Useful if the API is proxied through reverse proxy like nginx. Value needs to contain full URL with protocol scheme, e.g. https://braize.pajlada.com/chatterino")
	pflag.StringP("bind-address", "l", ":2558", "Address to which API will bind and start listening on")

	// Twitch
	pflag.String("twitch-login", "titlechange_bot", "Twitch login of the account on which bot will Log in to IRC")
	pflag.String("twitch-oauth", "", "OAuth token of the account on which bot will Log in to IRC. Should not have the \"oauth:\" part in the beginning.")
	pflag.String("twitch-client-id", "", "Twitch client ID")
	pflag.String("twitch-client-secret", "", "Twitch client secret")
	pflag.String("twitch-eventsub-secret", "", "Twitch EventSub secret used to create subscriptions and verify incoming notifications. Must be between 10 and 100 characters long")

	// Mongo 🥭
	pflag.String("mongo-username", "", "Username for the MongoDB user")
	pflag.String("mongo-password", "", "Password for the MongoDB user")
	pflag.String("mongo-port", "27017", "Port to which connection will try to connect. Note that you can only connect to localhost due to security concerns (use ssh port tunneling while testing/developing)")
	pflag.String("mongo-database-name", "tcb2", "Name of the database that should be used by the bot")
	pflag.String("mongo-auth-db", "admin", "Name of the authentication database, used as AuthSource while creating a new mongo.Connection. This should usually be left unchanged")

	pflag.Parse()
}

func New() (cfg *TCBConfig) {
	v := viper.New()

	v.BindPFlags(pflag.CommandLine)

	// figure out XDG_DATA_CONFIG to be compliant with the standard
	xdgConfigHome, exists := os.LookupEnv("XDG_CONFIG_HOME")
	if !exists || xdgConfigHome == "" {
		// on Windows, we use appdata since that's the closest equivalent
		if runtime.GOOS == "windows" {
			xdgConfigHome = "$APPDATA"
		} else {
			xdgConfigHome = filepath.Join("$HOME", ".config")
		}
	}

	// config paths to read from, in order of least importance
	var configPaths []string
	if runtime.GOOS != "windows" {
		configPaths = append(configPaths, filepath.Join("/etc", appName))
	}
	configPaths = append(configPaths, filepath.Join(xdgConfigHome, appName))
	configPaths = append(configPaths, ".")

	mergeConfig(v, configPaths)

	// Environment
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	v.UnmarshalExact(&cfg)

	//fmt.Printf("%# v\n", cfg) // uncomment for debugging purposes
	return
}
