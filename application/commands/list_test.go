package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dietrichm/admirer/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	t.Run("lists loved tracks for specified service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		tracks := []domain.Track{
			{
				Artist: "Awesome Artist",
				Name:   "Blam (Instrumental)",
			},
			{
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

		expected := `Awesome Artist - Blam (Instrumental)
Foo & Bar - Mr. Testy
`

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("returns error for unknown service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		expected := "unknown service"
		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(nil, errors.New(expected))

		output, err := executeList(serviceLoader, 5, "foobar")

		assert.EqualError(t, err, expected)
		assert.Empty(t, output)
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

		assert.Error(t, err)
		assert.Empty(t, output)
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

		assert.Error(t, err)
		assert.Empty(t, output)
	})
}

func executeList(serviceLoader domain.ServiceLoader, limit int, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := list(serviceLoader, limit, buffer, args)
	return buffer.String(), err
}
