package config

import (
	"errors"
	"testing"

	keyring_lib "github.com/99designs/keyring"
	gomock "github.com/golang/mock/gomock"
)

func TestKeyringConfig(t *testing.T) {
	t.Run("returns false when key is not set in keyring", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		keyring := NewMockKeyring(ctrl)
		item := keyring_lib.Item{}
		keyring.EXPECT().Get("foo").Return(item, errors.New("key error"))

		config := &keyringConfig{keyring}

		exists := config.IsSet("foo")

		if exists {
			t.Error("Key should not exist in keyring")
		}
	})

	t.Run("returns true when key is set in keyring", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		keyring := NewMockKeyring(ctrl)
		item := keyring_lib.Item{}
		keyring.EXPECT().Get("foo").Return(item, nil)

		config := &keyringConfig{keyring}

		exists := config.IsSet("foo")

		if !exists {
			t.Error("Key should exist in keyring")
		}
	})
}
