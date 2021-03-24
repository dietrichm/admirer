package spotify

import (
	"errors"
	"testing"

	mock_spotify "github.com/dietrichm/admirer/mock_services/spotify"
	"github.com/golang/mock/gomock"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
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

	t.Run("authenticates using authorization code", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		token := new(oauth2.Token)
		client := spotify.Client{}
		authenticator := mock_spotify.NewMockAuthenticator(ctrl)
		authenticator.EXPECT().Exchange("authcode").Return(token, nil)
		authenticator.EXPECT().NewClient(token).Return(client)

		service := &Spotify{authenticator: authenticator}

		err := service.Authenticate("authcode")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if service.client == nil {
			t.Error("Expected client to be set")
		}
	})

	t.Run("returns error for invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		authenticator := mock_spotify.NewMockAuthenticator(ctrl)
		authenticator.EXPECT().Exchange(gomock.Any()).Return(nil, errors.New("error"))

		service := &Spotify{authenticator: authenticator}

		err := service.Authenticate("authcode")

		if err == nil {
			t.Fatal("Expected an error")
		}
	})

	t.Run("returns username from client", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		user := &spotify.PrivateUser{
			User: spotify.User{
				DisplayName: "Joe",
			},
		}

		client := mock_spotify.NewMockClient(ctrl)
		client.EXPECT().CurrentUser().Return(user, nil)

		service := &Spotify{client: client}

		expected := "Joe"
		got, err := service.GetUsername()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}
