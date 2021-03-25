package lastfm

import (
	"errors"
	"testing"

	mock_lastfm "github.com/dietrichm/admirer/mock_services/lastfm"
	"github.com/golang/mock/gomock"
)

func TestLastfm(t *testing.T) {
	t.Run("creates authentication URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := mock_lastfm.NewMockAPI(ctrl)
		api.EXPECT().GetAuthRequestUrl("https://admirer.test").Return("https://service.test/auth")

		service := &Lastfm{api: api}

		got := service.CreateAuthURL()
		expected := "https://service.test/auth"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("authenticates using authorization code", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := mock_lastfm.NewMockAPI(ctrl)
		api.EXPECT().LoginWithToken("authcode")

		service := &Lastfm{api: api}

		err := service.Authenticate("authcode")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error for invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := mock_lastfm.NewMockAPI(ctrl)
		api.EXPECT().LoginWithToken(gomock.Any()).Return(errors.New("error"))

		service := &Lastfm{api: api}

		err := service.Authenticate("authcode")

		if err == nil {
			t.Fatal("Expected an error")
		}
	})
}
