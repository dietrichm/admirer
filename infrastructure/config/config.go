package config

import (
	"github.com/spf13/viper"
)

// Config is the interface for reading and writing configuration.
type Config interface {
	GetString(key string) string
}

func loadConfigFromFile(filename string) (Config, error) {
	config := viper.New()
	config.SetConfigFile(filename)
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}
