package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dietrichm/admirer/domain"
	"github.com/golang/mock/gomock"
)

func TestList(t *testing.T) {
	t.Run("lists loved tracks for specified service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		tracks := []domain.Track{
			domain.Track{
				Artist: "Awesome Artist",
				Name:   "Blam (Instrumental)",
			},
			domain.Track{
				Artist: "Foo & Bar",
				Name:   "Mr. Testy",
			},
		}

		service := domain.NewMockService(ctrl)
		service.EXPECT().Authenticated().Return(true)
		service.EXPECT().GetLovedTracks(5).Return(tracks, nil)
		service.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("foo").Return(service, nil)

		got, err := executeList(serviceLoader, 5, "foo")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := `Awesome Artist - Blam (Instrumental)
Foo & Bar - Mr. Testy
`
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error for unknown service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		expected := "unknown service"
		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(nil, errors.New(expected))

		output, err := executeList(serviceLoader, 5, "foobar")

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}

		if err == nil {
			t.Fatal("Expected an error")
		}

		got := err.Error()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error when service is not authenticated", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		service := domain.NewMockService(ctrl)
		service.EXPECT().Authenticated().Return(false)
		service.EXPECT().Name().AnyTimes().Return("Foo")
		service.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("foo").Return(service, nil)

		output, err := executeList(serviceLoader, 3, "foo")

		if err == nil {
			t.Error("Expected an error")
		}

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}
	})

	t.Run("returns error when failing to load loved tracks", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		service := domain.NewMockService(ctrl)
		service.EXPECT().Authenticated().Return(true)
		service.EXPECT().GetLovedTracks(3).Return(nil, errors.New("load error"))
		service.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("foo").Return(service, nil)

		output, err := executeList(serviceLoader, 3, "foo")

		if err == nil {
			t.Error("Expected an error")
		}

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}
	})
}

func executeList(serviceLoader domain.ServiceLoader, limit int, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := list(serviceLoader, limit, buffer, args)
	return buffer.String(), err
}
