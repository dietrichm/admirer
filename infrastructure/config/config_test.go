package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("reads YAML configuration file", func(t *testing.T) {
		file, err := createFile("foo: bar")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		defer os.Remove(file.Name())

		config, err := loadConfigFromFile(file.Name())

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if config == nil {
			t.Fatal("Expected config struct")
		}

		expected := "bar"
		got := config.GetString("foo")
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

func createFile(contents string) (*os.File, error) {
	file, err := ioutil.TempFile(os.TempDir(), "admirer-")
	if err != nil {
		return nil, err
	}

	if _, err := file.WriteString(contents); err != nil {
		return nil, err
	}

	return file, nil
}
