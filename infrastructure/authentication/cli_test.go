package authentication

import (
	"bytes"
	"strings"
	"testing"
)

func TestCliCallbackProvider(t *testing.T) {
	t.Run("returns code read from CLI iput", func(t *testing.T) {
		buffer := new(bytes.Buffer)
		buffer.WriteString("readString\n")
		writer := new(bytes.Buffer)

		provider := &cliCallbackProvider{buffer}

		got, err := provider.ReadCode("foo", writer)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "readString"
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		output := writer.String()
		expected = `Please provide "foo" parameter`

		if !strings.Contains(output, expected) {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("returns error when failing to read", func(t *testing.T) {
		buffer := new(bytes.Buffer)
		writer := new(bytes.Buffer)

		provider := &cliCallbackProvider{buffer}

		got, err := provider.ReadCode("foo", writer)

		if err == nil {
			t.Error("Expected an error")
		}

		if got != "" {
			t.Errorf("Unexpected result: %q", got)
		}
	})
}
