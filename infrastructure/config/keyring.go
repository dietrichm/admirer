//go:generate mockgen -source keyring.go -destination keyring_mock.go -package config

package config

import (
	"github.com/99designs/keyring"
)

// Keyring is our interface for a Keyring implementation.
type Keyring interface {
	Get(key string) (keyring.Item, error)
}

type keyringConfig struct {
	Keyring
}
