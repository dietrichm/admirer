package commands

import (
	"bytes"
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
		service.EXPECT().GetLovedTracks().Return(tracks)
		service.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("foo").Return(service, nil)

		got, err := executeList(serviceLoader, "foo")

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
}

func executeList(serviceLoader domain.ServiceLoader, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := list(serviceLoader, buffer, args)
	return buffer.String(), err
}
