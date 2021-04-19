//go:generate mockgen -source keyring.go -destination keyring_mock.go -package config

package config

import (
	"fmt"

	"github.com/99designs/keyring"
)

// Keyring is our interface for a Keyring implementation.
type Keyring interface {
	Get(key string) (keyring.Item, error)
	Set(item keyring.Item) error
}

type keyringConfig struct {
	Keyring
	prefix  string
	unsaved map[string]keyring.Item
}

func (k *keyringConfig) Get(key string) (keyring.Item, error) {
	if k.unsaved != nil {
		if item, exists := k.unsaved[key]; exists {
			return item, nil
		}
	}

	return k.Keyring.Get(key)
}

func (k *keyringConfig) IsSet(key string) bool {
	if _, err := k.Get(k.prefixed(key)); err != nil {
		return false
	}

	return true
}

func (k *keyringConfig) GetString(key string) string {
	item, err := k.Get(k.prefixed(key))
	if err != nil {
		return ""
	}

	return string(item.Data)
}

func (k *keyringConfig) prefixed(key string) string {
	return fmt.Sprintf("%s-%s", k.prefix, key)
}

func (k *keyringConfig) Set(key string, value interface{}) {
	item := keyring.Item{
		Key:  k.prefixed(key),
		Data: []byte(value.(string)),
	}

	if k.unsaved == nil {
		k.unsaved = map[string]keyring.Item{}
	}
	k.unsaved[key] = item
}

func (k *keyringConfig) Save() error {
	if k.unsaved != nil {
		for _, item := range k.unsaved {
			if err := k.Keyring.Set(item); err != nil {
				return err
			}
		}
		k.unsaved = map[string]keyring.Item{}
	}

	return nil
}

type keyringLoader struct{}

func (k keyringLoader) Load(name string) (Config, error) {
	config := keyring.Config{
		ServiceName: "admirer",
	}

	return k.open(name, config)
}

func (k keyringLoader) open(name string, config keyring.Config) (Config, error) {
	ring, err := keyring.Open(config)
	if err != nil {
		return nil, fmt.Errorf("failed opening keyring: %w", err)
	}

	return &keyringConfig{
		Keyring: ring,
		prefix:  name,
	}, nil
}
