package commands

import (
	"bytes"
	"errors"
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

	t.Run("returns error when failing to mark track as loved", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		tracks := []domain.Track{
			domain.Track{
				Artist: "Awesome Artist",
				Name:   "Blam (Instrumental)",
			},
		}

		sourceService := domain.NewMockService(ctrl)
		sourceService.EXPECT().GetLovedTracks(gomock.Any()).Return(tracks, nil)
		sourceService.EXPECT().Close()

		targetService := domain.NewMockService(ctrl)
		targetService.EXPECT().LoveTrack(gomock.Any()).Return(errors.New("api error"))
		targetService.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(targetService, nil)

		output, err := executeSync(serviceLoader, "source", "target")

		if err == nil {
			t.Error("Expected an error")
		}

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}
	})

	t.Run("returns error when failing to read loved tracks", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		sourceService := domain.NewMockService(ctrl)
		sourceService.EXPECT().GetLovedTracks(gomock.Any()).Return(nil, errors.New("read error"))
		sourceService.EXPECT().Close()

		targetService := domain.NewMockService(ctrl)
		targetService.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(targetService, nil)

		output, err := executeSync(serviceLoader, "source", "target")

		if err == nil {
			t.Error("Expected an error")
		}

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}
	})
}

func executeSync(serviceLoader domain.ServiceLoader, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := sync(serviceLoader, buffer, args)
	return buffer.String(), err
}
