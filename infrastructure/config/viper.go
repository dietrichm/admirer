package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	permissions          os.FileMode = 0600
	permissionsString    string      = "-rw-------"
	directoryPermissions os.FileMode = 0700
)

type viperConfig struct {
	*viper.Viper
}

func (v *viperConfig) Save() error {
	return v.WriteConfig()
}

type viperLoader struct{}

func (v viperLoader) Load(name string) (Config, error) {
	filename := filepath.Join(os.Getenv("HOME"), ".config", "admirer", name)
	return v.loadFromFile(filename)
}

func (v viperLoader) loadFromFile(filename string) (Config, error) {
	config := viper.New()
	config.SetConfigFile(filename)
	config.SetConfigType("yaml")
	config.SetConfigPermissions(permissions)

	if err := config.ReadInConfig(); err != nil {
		switch readError := err.(type) {
		case *fs.PathError:
			directory := filepath.Dir(filename)
			if err := os.MkdirAll(directory, directoryPermissions); err != nil {
				return nil, fmt.Errorf("failed creating configuration directory: %w", err)
			}
			if err := config.WriteConfig(); err != nil {
				return nil, fmt.Errorf("failed creating configuration file: %w", err)
			}
		default:
			return nil, fmt.Errorf("failed reading configuration file: %w", readError)
		}
	}

	if err := v.checkPermissions(filename); err != nil {
		return nil, err
	}

	return &viperConfig{config}, nil
}

func (v viperLoader) checkPermissions(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed opening file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed opening file: %w", err)
	}

	actualPermissions := stat.Mode().Perm()

	if actualPermissions.String() != permissionsString {
		return fmt.Errorf("wrong permissions on %q: got %q, want %q", filename, actualPermissions, permissionsString)
	}

	return nil
}
