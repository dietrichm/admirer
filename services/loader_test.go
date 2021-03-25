package services

import (
	"testing"

	"github.com/dietrichm/admirer/services/spotify"
)

func TestMapServiceLoader(t *testing.T) {
	t.Run("returns service when loader exists", func(t *testing.T) {
		service := &spotify.Spotify{}

		serviceLoader := MapServiceLoader{
			"foo": func() (Service, error) {
				return service, nil
			},
			"bar": func() (Service, error) {
				return nil, nil
			},
		}

		expected := service
		got, err := serviceLoader.ForName("foo")

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
