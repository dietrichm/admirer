package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dietrichm/admirer/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
		sourceService.EXPECT().Authenticated().Return(true)
		sourceService.EXPECT().GetLovedTracks(5).Return(tracks, nil)
		sourceService.EXPECT().Close()

		targetService := domain.NewMockService(ctrl)
		targetService.EXPECT().Authenticated().Return(true)
		targetService.EXPECT().LoveTrack(trackOne).Return(nil)
		targetService.EXPECT().LoveTrack(trackTwo).Return(nil)
		targetService.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(targetService, nil)

		got, err := executeSync(serviceLoader, 5, "source", "target")

		expected := `Synced: Awesome Artist - Blam (Instrumental)
Synced: Foo & Bar - Mr. Testy
`

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("returns error when failing to mark track as loved", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		tracks := []domain.Track{
			{
				Artist: "Awesome Artist",
				Name:   "Blam (Instrumental)",
			},
		}

		sourceService := domain.NewMockService(ctrl)
		sourceService.EXPECT().Authenticated().Return(true)
		sourceService.EXPECT().GetLovedTracks(gomock.Any()).Return(tracks, nil)
		sourceService.EXPECT().Close()

		targetService := domain.NewMockService(ctrl)
		targetService.EXPECT().Authenticated().Return(true)
		targetService.EXPECT().LoveTrack(gomock.Any()).Return(errors.New("api error"))
		targetService.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(targetService, nil)

		output, err := executeSync(serviceLoader, 10, "source", "target")

		assert.Error(t, err)
		assert.Empty(t, output)
	})

	t.Run("returns error when failing to read loved tracks", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		sourceService := domain.NewMockService(ctrl)
		sourceService.EXPECT().Authenticated().Return(true)
		sourceService.EXPECT().GetLovedTracks(gomock.Any()).Return(nil, errors.New("read error"))
		sourceService.EXPECT().Close()

		targetService := domain.NewMockService(ctrl)
		targetService.EXPECT().Authenticated().Return(true)
		targetService.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(targetService, nil)

		output, err := executeSync(serviceLoader, 10, "source", "target")

		assert.Error(t, err)
		assert.Empty(t, output)
	})

	t.Run("returns error when failing to load source service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(nil, errors.New("service error"))

		output, err := executeSync(serviceLoader, 10, "source", "target")

		assert.Error(t, err)
		assert.Empty(t, output)
	})

	t.Run("returns error when failing to load target service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		sourceService := domain.NewMockService(ctrl)
		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(nil, errors.New("service error"))

		output, err := executeSync(serviceLoader, 10, "source", "target")

		assert.Error(t, err)
		assert.Empty(t, output)
	})

	t.Run("returns error when source service is not authenticated", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		sourceService := domain.NewMockService(ctrl)
		sourceService.EXPECT().Authenticated().Return(false)
		sourceService.EXPECT().Name().AnyTimes().Return("Source")
		sourceService.EXPECT().Close()

		targetService := domain.NewMockService(ctrl)
		targetService.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(targetService, nil)

		output, err := executeSync(serviceLoader, 10, "source", "target")

		assert.Error(t, err)
		assert.Empty(t, output)
	})

	t.Run("returns error when target service is not authenticated", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		sourceService := domain.NewMockService(ctrl)
		sourceService.EXPECT().Authenticated().Return(true)
		sourceService.EXPECT().Close()

		targetService := domain.NewMockService(ctrl)
		targetService.EXPECT().Authenticated().Return(false)
		targetService.EXPECT().Name().AnyTimes().Return("Target")
		targetService.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("source").Return(sourceService, nil)
		serviceLoader.EXPECT().ForName("target").Return(targetService, nil)

		output, err := executeSync(serviceLoader, 10, "source", "target")

		assert.Error(t, err)
		assert.Empty(t, output)
	})
}

func executeSync(serviceLoader domain.ServiceLoader, limit int, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := sync(serviceLoader, limit, buffer, args)
	return buffer.String(), err
}
