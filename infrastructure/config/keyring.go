//go:generate mockgen -source keyring.go -destination keyring_mock.go -package config

package config

import (
	"fmt"

	"github.com/99designs/keyring"
)

// Keyring is our interface for a Keyring implementation.
type Keyring interface {
	Get(key string) (keyring.Item, error)
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

func (k *keyringConfig) prefixed(key string) string {
	return fmt.Sprintf("%s-%s", k.prefix, key)
}
