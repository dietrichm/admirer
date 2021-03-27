package config

import (
	"io/fs"
	"os"

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
		switch readError := err.(type) {
		case *fs.PathError:
			os.Create(filename)
		default:
			return nil, readError
		}
	}

	return config, nil
}
