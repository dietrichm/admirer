package commands

import (
	"bytes"
	"testing"

	"github.com/dietrichm/admirer/domain"
	"github.com/golang/mock/gomock"
)

func TestStatus(t *testing.T) {
	t.Run("returns status for each service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		fooService := domain.NewMockService(ctrl)
		fooService.EXPECT().Name().Return("Foo")
		fooService.EXPECT().GetUsername().Return("user303", nil)

		barService := domain.NewMockService(ctrl)
		barService.EXPECT().Name().Return("Bar")
		barService.EXPECT().GetUsername().Return("user808", nil)

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().Names().Return([]string{"foo", "bar"})
		serviceLoader.EXPECT().ForName("foo").Return(fooService, nil)
		serviceLoader.EXPECT().ForName("bar").Return(barService, nil)

		expected := `Foo
	Authenticated as user303
Bar
	Authenticated as user808
`
		got, err := executeStatus(serviceLoader)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

func executeStatus(serviceLoader domain.ServiceLoader) (string, error) {
	buffer := new(bytes.Buffer)
	err := status(serviceLoader, buffer)
	return buffer.String(), err
}
