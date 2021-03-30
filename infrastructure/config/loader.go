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

type viperLoader struct{}

// Load Config from file system.
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
				return nil, err
			}
			if err := config.WriteConfig(); err != nil {
				return nil, err
			}
		default:
			return nil, readError
		}
	}

	if err := v.checkPermissions(filename); err != nil {
		return nil, err
	}

	return config, nil
}

func (v viperLoader) checkPermissions(filename string) error {
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
