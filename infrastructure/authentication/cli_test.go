package authentication

import (
	"testing"
)

func TestCliCallbackProvider(t *testing.T) {
	t.Run("returns code read from CLI iput", func(t *testing.T) {
		provider := new(cliCallbackProvider)

		got, err := provider.ReadCode("foo")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := ""
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}
