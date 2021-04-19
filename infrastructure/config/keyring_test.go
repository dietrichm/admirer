package config

import (
	"errors"
	"os"
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

		config := &keyringConfig{
			Keyring: keyring,
			prefix:  "prefix",
		}

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

		config := &keyringConfig{
			Keyring: keyring,
			prefix:  "prefix",
		}

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

		config := &keyringConfig{
			Keyring: keyring,
			prefix:  "prefix",
		}

		config.Set("foo", "bar")
		err := config.Save()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error when failing to save key-value pair", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		keyring := NewMockKeyring(ctrl)
		keyring.EXPECT().Set(gomock.Any()).Return(errors.New("write error"))

		config := &keyringConfig{
			Keyring: keyring,
			prefix:  "prefix",
		}

		config.Set("foo", "bar")
		err := config.Save()

		if err == nil {
			t.Fatal("Expected an error")
		}
	})
}

func TestKeyringLoader(t *testing.T) {
	t.Run("opens keyring config with settings", func(t *testing.T) {
		loader := &keyringLoader{}

		keyringConfig := keyring_lib.Config{
			AllowedBackends: []keyring_lib.BackendType{keyring_lib.FileBackend},
			FileDir:         os.TempDir(),
		}
		config, err := loader.open("name", keyringConfig)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if config == nil {
			t.Error("Expected config instance")
		}
	})

	t.Run("returns error when unable to open keyring", func(t *testing.T) {
		loader := &keyringLoader{}

		keyringConfig := keyring_lib.Config{
			AllowedBackends: []keyring_lib.BackendType{},
		}
		config, err := loader.open("name", keyringConfig)

		if err == nil {
			t.Error("Expected an error")
		}

		if config != nil {
			t.Errorf("Unexpected config instance: %v", config)
		}
	})
}
