package spotify

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/config"
	"github.com/golang/mock/gomock"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func TestSpotify(t *testing.T) {
	t.Run("returns whether service is authenticated", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		service := &Spotify{}

		if service.Authenticated() {
			t.Error("expected not to be authenticated")
		}

		service = &Spotify{client: NewMockClient(ctrl)}

		if !service.Authenticated() {
			t.Error("expected to be authenticated")
		}
	})

	t.Run("creates authentication URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		redirectOption := oauth2.SetAuthURLParam("redirect_uri", "https://admirer.test/foo")
		authenticator := NewMockAuthenticator(ctrl)
		authenticator.EXPECT().AuthURLWithOpts(gomock.Any(), redirectOption).Return("https://service.test/auth")

		service := &Spotify{authenticator: authenticator}

		expected := "https://service.test/auth"
		got := service.CreateAuthURL("https://admirer.test/foo")

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("authenticates using authorization code", func(t *testing.T) {
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

		service := &Spotify{
			authenticator: authenticator,
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

	t.Run("authenticates from token in secrets", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		now := time.Now()
		token := &oauth2.Token{
			TokenType:    "myTokenType",
			AccessToken:  "myAccessToken",
			Expiry:       now.Truncate(time.Second),
			RefreshToken: "myRefreshToken",
		}
		client := spotify.Client{}

		authenticator := NewMockAuthenticator(ctrl)
		authenticator.EXPECT().NewClient(&tokenMatcher{token}).Return(client)

		secrets := config.NewMockConfig(ctrl)
		secrets.EXPECT().IsSet("token_type").Return(true)
		secrets.EXPECT().GetString("expiry").Return(now.Format(time.RFC3339))
		secrets.EXPECT().GetString("token_type").Return("myTokenType")
		secrets.EXPECT().GetString("access_token").Return("myAccessToken")
		secrets.EXPECT().GetString("refresh_token").Return("myRefreshToken")

		service := &Spotify{
			authenticator: authenticator,
			secrets:       secrets,
		}

		service.authenticateFromSecrets(secrets)
	})

	t.Run("returns loved tracks for current user", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		result := &spotify.SavedTrackPage{
			Tracks: []spotify.SavedTrack{
				spotify.SavedTrack{
					FullTrack: spotify.FullTrack{
						SimpleTrack: spotify.SimpleTrack{
							Artists: []spotify.SimpleArtist{
								spotify.SimpleArtist{
									Name: "Awesome Artist",
								},
							},
							Name: "Blam (Instrumental)",
						},
					},
				},
				spotify.SavedTrack{
					FullTrack: spotify.FullTrack{
						SimpleTrack: spotify.SimpleTrack{
							Artists: []spotify.SimpleArtist{
								spotify.SimpleArtist{
									Name: "Foo & Bar",
								},
							},
							Name: "Mr. Testy",
						},
					},
				},
			},
		}

		client := NewMockClient(ctrl)
		limit := 5
		options := &spotify.Options{
			Limit: &limit,
		}
		client.EXPECT().CurrentUsersTracksOpt(options).Return(result, nil)

		service := &Spotify{
			client: client,
		}

		expected := []domain.Track{
			domain.Track{
				Artist: "Awesome Artist",
				Name:   "Blam (Instrumental)",
			},
			domain.Track{
				Artist: "Foo & Bar",
				Name:   "Mr. Testy",
			},
		}
		got, err := service.GetLovedTracks(5)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		for index, expectedTrack := range expected {
			if got[index] != expectedTrack {
				t.Errorf("expected %q, got %q", expectedTrack, got[index])
			}
		}
	})

	t.Run("returns error when failing to read loved tracks", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		client := NewMockClient(ctrl)
		client.EXPECT().CurrentUsersTracksOpt(gomock.Any()).Return(nil, errors.New("read error"))

		service := &Spotify{
			client: client,
		}

		got, err := service.GetLovedTracks(5)

		if err == nil {
			t.Error("Expected an error")
		}

		if len(got) > 0 {
			t.Errorf("Unexpected return value: %v", got)
		}
	})

	t.Run("marks track as loved", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		limit := 1
		options := &spotify.Options{
			Limit: &limit,
		}
		result := &spotify.SearchResult{
			Tracks: &spotify.FullTrackPage{
				Tracks: []spotify.FullTrack{
					spotify.FullTrack{
						SimpleTrack: spotify.SimpleTrack{
							ID: "trackID",
						},
					},
				},
			},
		}

		client := NewMockClient(ctrl)
		client.EXPECT().SearchOpt(`artist:"Foo & Bar The Famous Two" track:"Mr. Testy - 12 Version"`, gomock.Any(), options).Return(result, nil)
		client.EXPECT().AddTracksToLibrary([]spotify.ID{"trackID"})

		service := &Spotify{
			client: client,
		}

		track := domain.Track{
			Artist: `Foo & Bar "The Famous Two"`,
			Name:   `Mr. Testy - 12" Version`,
		}

		err := service.LoveTrack(track)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error when failing to search track", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		client := NewMockClient(ctrl)
		client.EXPECT().SearchOpt(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("read error"))

		service := &Spotify{
			client: client,
		}

		track := domain.Track{
			Artist: "Foo & Bar",
			Name:   "Mr. Testy",
		}

		err := service.LoveTrack(track)

		if err == nil {
			t.Error("Expected an error")
		}
	})

	t.Run("skip marking track as loved when no track is found", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		result := &spotify.SearchResult{
			Tracks: &spotify.FullTrackPage{
				Tracks: []spotify.FullTrack{},
			},
		}

		client := NewMockClient(ctrl)
		client.EXPECT().SearchOpt(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)

		service := &Spotify{
			client: client,
		}

		track := domain.Track{
			Artist: "Foo & Bar",
			Name:   "Mr. Testy",
		}

		err := service.LoveTrack(track)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error when failing to add track to library", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		result := &spotify.SearchResult{
			Tracks: &spotify.FullTrackPage{
				Tracks: []spotify.FullTrack{
					spotify.FullTrack{
						SimpleTrack: spotify.SimpleTrack{
							ID: "trackID",
						},
					},
				},
			},
		}

		client := NewMockClient(ctrl)
		client.EXPECT().SearchOpt(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)
		client.EXPECT().AddTracksToLibrary(gomock.Any()).Return(errors.New("api error"))

		service := &Spotify{
			client: client,
		}

		track := domain.Track{
			Artist: "Foo & Bar",
			Name:   "Mr. Testy",
		}

		err := service.LoveTrack(track)

		if err == nil {
			t.Error("Expected an error")
		}
	})

	t.Run("new token is persisted when closing service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		now := time.Now()
		token := &oauth2.Token{
			TokenType:    "myTokenType",
			AccessToken:  "myAccessToken",
			Expiry:       now,
			RefreshToken: "myRefreshToken",
		}

		client := NewMockClient(ctrl)
		client.EXPECT().Token().Return(token, nil)

		secrets := config.NewMockConfig(ctrl)
		gomock.InOrder(
			secrets.EXPECT().Set("token_type", "myTokenType"),
			secrets.EXPECT().Set("access_token", "myAccessToken"),
			secrets.EXPECT().Set("expiry", now.Format(time.RFC3339)),
			secrets.EXPECT().Set("refresh_token", "myRefreshToken"),
			secrets.EXPECT().Save(),
		)

		service := &Spotify{
			client:  client,
			secrets: secrets,
		}

		err := service.Close()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("skip persisting token on close when not authenticated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		secrets := config.NewMockConfig(ctrl)

		service := &Spotify{
			secrets: secrets,
		}

		err := service.Close()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error when unable to read token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		client := NewMockClient(ctrl)
		client.EXPECT().Token().Return(nil, errors.New("token error"))

		service := &Spotify{
			client: client,
		}

		err := service.Close()

		if err == nil {
			t.Fatal("Expected an error")
		}
	})

	t.Run("returns error when unable to save secrets", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		now := time.Now()
		token := &oauth2.Token{
			TokenType:    "myTokenType",
			AccessToken:  "myAccessToken",
			Expiry:       now,
			RefreshToken: "myRefreshToken",
		}

		client := NewMockClient(ctrl)
		client.EXPECT().Token().Return(token, nil)

		secrets := config.NewMockConfig(ctrl)
		secrets.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
		secrets.EXPECT().Save().Return(errors.New("write error"))

		service := &Spotify{
			client:  client,
			secrets: secrets,
		}

		err := service.Close()

		if err == nil {
			t.Fatal("Expected an error")
		}
	})
}

func TestNewSpotify(t *testing.T) {
	t.Run("creates instance when environment is configured", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		secrets := config.NewMockConfig(ctrl)
		secrets.EXPECT().IsSet("token_type").Return(false)

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

type tokenMatcher struct {
	token *oauth2.Token
}

func (t tokenMatcher) Matches(x interface{}) bool {
	tokenWithoutExpiry := &oauth2.Token{
		TokenType:    t.token.TokenType,
		AccessToken:  t.token.AccessToken,
		RefreshToken: t.token.RefreshToken,
	}

	gotToken, ok := x.(*oauth2.Token)
	if !ok {
		return false
	}

	gotTokenWithoutExpiry := &oauth2.Token{
		TokenType:    gotToken.TokenType,
		AccessToken:  gotToken.AccessToken,
		RefreshToken: gotToken.RefreshToken,
	}

	if !reflect.DeepEqual(tokenWithoutExpiry, gotTokenWithoutExpiry) {
		return false
	}

	expiry := t.token.Expiry.Truncate(time.Second)
	gotExpiry := gotToken.Expiry.Truncate(time.Second)

	if !expiry.Equal(gotExpiry) {
		return false
	}

	return true
}

func (t tokenMatcher) String() string {
	return fmt.Sprintf("is equal to %v", t.token)
}
