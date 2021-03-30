package spotify

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/dietrichm/admirer/infrastructure/config"
	"github.com/golang/mock/gomock"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func TestSpotify(t *testing.T) {
	t.Run("creates authentication URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		authenticator := NewMockAuthenticator(ctrl)
		authenticator.EXPECT().AuthURL("").Return("https://service.test/auth")

		service := &Spotify{authenticator: authenticator}

		expected := "https://service.test/auth"
		got := service.CreateAuthURL()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("authenticates using authorization code and saves oauth token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		now := time.Now()
		token := &oauth2.Token{
			TokenType:    "myTokenType",
			AccessToken:  "myAccessToken",
			Expiry:       now,
			RefreshToken: "myRefreshToken",
		}

		client := spotify.Client{}
		authenticator := NewMockAuthenticator(ctrl)
		authenticator.EXPECT().Exchange("authcode").Return(token, nil)
		authenticator.EXPECT().NewClient(token).Return(client)

		secrets := config.NewMockConfig(ctrl)
		gomock.InOrder(
			secrets.EXPECT().Set("token_type", "myTokenType"),
			secrets.EXPECT().Set("access_token", "myAccessToken"),
			secrets.EXPECT().Set("expiry", now),
			secrets.EXPECT().Set("refresh_token", "myRefreshToken"),
			secrets.EXPECT().WriteConfig(),
		)

		service := &Spotify{
			authenticator: authenticator,
			secrets:       secrets,
		}

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

		authenticator := NewMockAuthenticator(ctrl)
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

		client := NewMockClient(ctrl)
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

	t.Run("returns error when failing to read username", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		client := NewMockClient(ctrl)
		client.EXPECT().CurrentUser().Return(nil, errors.New("error"))

		service := &Spotify{client: client}

		got, err := service.GetUsername()
		expected := ""

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if err == nil {
			t.Fatal("Expected an error")
		}
	})
}

func TestNewSpotify(t *testing.T) {
	t.Run("creates instance when environment is configured", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		secrets := config.NewMockConfig(ctrl)

		os.Setenv("SPOTIFY_CLIENT_ID", "client_id")
		os.Setenv("SPOTIFY_CLIENT_SECRET", "client_secret")

		service, err := NewSpotify(secrets)

		if service == nil {
			t.Error("Expected an instance")
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error when environment is not configured", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		secrets := config.NewMockConfig(ctrl)

		os.Unsetenv("SPOTIFY_CLIENT_ID")
		os.Unsetenv("SPOTIFY_CLIENT_SECRET")

		service, err := NewSpotify(secrets)

		if service != nil {
			t.Errorf("Unexpected instance: %v", service)
		}

		if err == nil {
			t.Error("Expected an error")
		}
	})
}
