package spotify

import (
	"testing"

	mock_spotify "github.com/dietrichm/admirer/mock_services/spotify"
	"github.com/golang/mock/gomock"
)

func TestSpotify(t *testing.T) {
	t.Run("creates authentication URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		authenticator := mock_spotify.NewMockAuthenticator(ctrl)
		authenticator.EXPECT().AuthURL("").Return("https://service.test/auth")

		service := &Spotify{authenticator: authenticator}

		expected := "https://service.test/auth"
		got := service.CreateAuthURL()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}
