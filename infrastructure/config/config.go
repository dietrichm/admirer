//go:generate mockgen -source config.go -destination config_mock.go -package config

package config

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config is the interface for reading and writing configuration.
type Config interface {
	GetString(key string) string
	Set(key string, value interface{})
	WriteConfig() error
}

// LoadConfig creates a Config struct for reading and writing configuration.
func LoadConfig(name string) (Config, error) {
	filename := filepath.Join(os.Getenv("HOME"), ".config", "admirer", name)
	return loadConfigFromFile(filename)
}

func loadConfigFromFile(filename string) (Config, error) {
	config := viper.New()
	config.SetConfigFile(filename)
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		switch readError := err.(type) {
		case *fs.PathError:
			if _, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600); err != nil {
				return nil, err
			}
		default:
			return nil, readError
		}
	}

	return config, nil
}
