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
	prefix string
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

	k.Keyring.Set(item)
}

func (k *keyringConfig) Save() error {
	return nil
}
