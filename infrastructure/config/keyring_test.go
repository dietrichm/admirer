package config

import (
	"errors"
	"testing"

	keyring_lib "github.com/99designs/keyring"
	gomock "github.com/golang/mock/gomock"
)

func TestKeyringConfig(t *testing.T) {
	t.Run("returns whether key is set in keyring", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		keyring := NewMockKeyring(ctrl)
		item := keyring_lib.Item{}
		keyring.EXPECT().Get("prefix-foo").Return(item, errors.New("key error"))
		keyring.EXPECT().Get("prefix-bar").Return(item, nil)

		config := &keyringConfig{keyring, "prefix"}

		exists := config.IsSet("foo")

		if exists {
			t.Error("Key should not exist in keyring")
		}

		exists = config.IsSet("bar")

		if !exists {
			t.Error("Key should exist in keyring")
		}
	})

	t.Run("returns string value for key", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		keyring := NewMockKeyring(ctrl)
		item := keyring_lib.Item{
			Data: []byte("something"),
		}
		keyring.EXPECT().Get("prefix-foo").Return(item, errors.New("key error"))
		keyring.EXPECT().Get("prefix-bar").Return(item, nil)

		config := &keyringConfig{keyring, "prefix"}

		got := config.GetString("foo")
		expected := ""

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		got = config.GetString("bar")
		expected = "something"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("saves key-value pair in keyring", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		keyring := NewMockKeyring(ctrl)
		item := keyring_lib.Item{
			Key:  "prefix-foo",
			Data: []byte("bar"),
		}
		keyring.EXPECT().Set(item).Return(nil)

		config := &keyringConfig{keyring, "prefix"}

		config.Set("foo", "bar")
		err := config.Save()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
