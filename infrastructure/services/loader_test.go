package services

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMapServiceLoader(t *testing.T) {
	t.Run("returns service when loader exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := NewMockService(ctrl)

		serviceLoader := MapServiceLoader{
			"foo": func() (Service, error) {
				return service, nil
			},
			"bar": func() (Service, error) {
				return nil, nil
			},
		}

		got, err := serviceLoader.ForName("foo")

		if got != service {
			t.Errorf("expected %v, got %v", service, got)
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error when loader does not exist", func(t *testing.T) {
		serviceLoader := MapServiceLoader{
			"foo": func() (Service, error) {
				return nil, nil
			},
		}

		service, err := serviceLoader.ForName("bar")

		if service != nil {
			t.Errorf("Unexpected service instance: %v", service)
		}

		if err == nil {
			t.Fatal("Expected an error")
		}

		expected := `unknown service "bar"`
		got := err.Error()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error when loader yields error", func(t *testing.T) {
		serviceError := errors.New("service error")
		serviceLoader := MapServiceLoader{
			"foo": func() (Service, error) {
				return nil, serviceError
			},
		}

		service, err := serviceLoader.ForName("foo")

		if service != nil {
			t.Errorf("Unexpected service instance: %v", service)
		}

		if err == nil {
			t.Fatal("Expected an error")
		}

		if err != serviceError {
			t.Errorf("expected %q, got %q", serviceError, err)
		}
	})
}
