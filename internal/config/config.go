package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

func New() *TCBConfig {
	// Read the config file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	// Parse the config file (initialize the object with default values)
	cfg := &TCBConfig{
		CommandPrefix:     "!",
		BindAddress:       "localhost:2558",
		MongoPort:         "27017",
		MongoDatabaseName: "tcb2",
		MongoAuthDB:       "admin",
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal config file: %w", err))
	}

	return cfg
}
