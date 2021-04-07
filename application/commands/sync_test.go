package commands

import (
	"bytes"
	"testing"

	"github.com/dietrichm/admirer/domain"
	"github.com/golang/mock/gomock"
)

func TestSync(t *testing.T) {
	t.Run("loves tracks on target service returned from source service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		trackOne := domain.Track{
			Artist: "Awesome Artist",
			Name:   "Blam (Instrumental)",
		}
		trackTwo := domain.Track{
			Artist: "Foo & Bar",
			Name:   "Mr. Testy",
		}
		tracks := []domain.Track{trackOne, trackTwo}

		sourceService := domain.NewMockService(ctrl)
		sourceService.EXPECT().GetLovedTracks(10).Return(tracks, nil)
		sourceService.EXPECT().Close()

		targetService := domain.NewMockService(ctrl)
		targetService.EXPECT().LoveTrack(trackOne).Return(nil)
		targetService.EXPECT().LoveTrack(trackTwo).Return(nil)
		targetService.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(targetService, nil)

		got, err := executeSync(serviceLoader, "source", "target")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := `Synced: Awesome Artist - Blam (Instrumental)
Synced: Foo & Bar - Mr. Testy
`

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

func executeSync(serviceLoader domain.ServiceLoader, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := sync(serviceLoader, buffer, args)
	return buffer.String(), err
}
