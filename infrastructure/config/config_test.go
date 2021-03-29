package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

	t.Run("creates non existing configuration file with correct permissions", func(t *testing.T) {
		file, err := createFile("")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		os.Remove(file.Name())
		defer os.Remove(file.Name())

		config, err := loadConfigFromFile(file.Name())

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if config == nil {
			t.Fatal("Expected config struct")
		}

		file, _ = os.Open(file.Name())
		stat, _ := file.Stat()
		expected := "-rw-------"
		got := stat.Mode().Perm().String()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("creates non existing containing directory with correct permissions", func(t *testing.T) {
		file, err := createFile("")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		os.Remove(file.Name())

		directory := filepath.Join(filepath.Dir(file.Name()), "admirer-test")
		directoryFile := filepath.Join(directory, filepath.Base(file.Name()))
		defer os.RemoveAll(directory)

		config, err := loadConfigFromFile(directoryFile)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if config == nil {
			t.Fatal("Expected config struct")
		}

		file, _ = os.Open(directory)
		stat, _ := file.Stat()
		expected := "-rwx------"
		got := stat.Mode().Perm().String()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error for invalid filename", func(t *testing.T) {
		config, err := loadConfigFromFile("./")

		if err == nil {
			t.Fatal("Expected an error")
		}

		if config != nil {
			t.Errorf("Unexpected config struct: %v", config)
		}
	})

	t.Run("returns error for other read issues", func(t *testing.T) {
		file, err := createFile("$$$")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		defer os.Remove(file.Name())

		config, err := loadConfigFromFile(file.Name())

		if err == nil {
			t.Fatal("Expected an error")
		}

		if config != nil {
			t.Errorf("Unexpected config struct: %v", config)
		}
	})

	t.Run("returns error when file has incorrect permissions", func(t *testing.T) {
		file, err := createFile("")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		file.Chmod(0666)
		defer os.Remove(file.Name())

		config, err := loadConfigFromFile(file.Name())

		if err == nil {
			t.Fatal("Expected an error")
		}

		expected := "wrong permissions on "
		got := err.Error()

		if !strings.Contains(got, expected) {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if config != nil {
			t.Errorf("Unexpected config struct: %v", config)
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
