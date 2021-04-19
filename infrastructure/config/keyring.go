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
	err    error
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

	if err := k.Keyring.Set(item); err != nil {
		k.err = err
	}
}

func (k *keyringConfig) Save() error {
	return k.err
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
