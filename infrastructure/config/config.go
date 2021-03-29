//go:generate mockgen -source config.go -destination config_mock.go -package config

package config

import (
	"fmt"
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

// Loader loads Config by name.
type Loader interface {
	Load(name string) (Config, error)
}

type viperLoader struct{}

// Load Config from file system.
func (v viperLoader) Load(name string) (Config, error) {
	filename := filepath.Join(os.Getenv("HOME"), ".config", "admirer", name)
	return loadConfigFromFile(filename)
}

const (
	permissions          os.FileMode = 0600
	permissionsString    string      = "-rw-------"
	directoryPermissions os.FileMode = 0700
)

// LoadConfig creates a Config struct for reading and writing configuration.
func LoadConfig(name string) (Config, error) {
	filename := filepath.Join(os.Getenv("HOME"), ".config", "admirer", name)
	return loadConfigFromFile(filename)
}

func loadConfigFromFile(filename string) (Config, error) {
	config := viper.New()
	config.SetConfigFile(filename)
	config.SetConfigType("yaml")
	config.SetConfigPermissions(permissions)

	if err := config.ReadInConfig(); err != nil {
		switch readError := err.(type) {
		case *fs.PathError:
			directory := filepath.Dir(filename)
			if err := os.MkdirAll(directory, directoryPermissions); err != nil {
				return nil, err
			}
			if err := config.WriteConfig(); err != nil {
				return nil, err
			}
		default:
			return nil, readError
		}
	}

	if err := checkPermissions(filename); err != nil {
		return nil, err
	}

	return config, nil
}

func checkPermissions(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	actualPermissions := stat.Mode().Perm()

	if actualPermissions.String() != permissionsString {
		return fmt.Errorf("wrong permissions on %q: got %q, want %q", filename, actualPermissions, permissionsString)
	}

	return nil
}
