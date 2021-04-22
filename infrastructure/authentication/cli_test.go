package authentication

import (
	"bytes"
	"testing"
)

func TestCliCallbackProvider(t *testing.T) {
	t.Run("returns code read from CLI iput", func(t *testing.T) {
		buffer := new(bytes.Buffer)
		buffer.WriteString("readString\n")

		provider := &cliCallbackProvider{buffer}

		got, err := provider.ReadCode("foo")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "readString"
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error when failing to read", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		provider := &cliCallbackProvider{buffer}

		got, err := provider.ReadCode("foo")

		if err == nil {
			t.Error("Expected an error")
		}

		if got != "" {
			t.Errorf("Unexpected result: %q", got)
		}
	})
}
