package lastfm

import (
	"errors"
	"os"
	"testing"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/config"
	"github.com/golang/mock/gomock"
	"github.com/shkh/lastfm-go/lastfm"
)

func TestLastfm(t *testing.T) {
	t.Run("returns whether service is authenticated", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := NewMockAPI(ctrl)
		api.EXPECT().GetSessionKey().Return("")

		service := &Lastfm{api: api}

		if service.Authenticated() {
			t.Error("expected not to be authenticated")
		}

		api = NewMockAPI(ctrl)
		api.EXPECT().GetSessionKey().Return("sessionKey")

		service = &Lastfm{api: api}

		if !service.Authenticated() {
			t.Error("expected to be authenticated")
		}
	})

	t.Run("creates authentication URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := NewMockAPI(ctrl)
		api.EXPECT().GetAuthRequestUrl("https://admirer.test").Return("https://service.test/auth")

		service := &Lastfm{api: api}

		got := service.CreateAuthURL()
		expected := "https://service.test/auth"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("authenticates using authorization code and saves session key", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := NewMockAPI(ctrl)
		api.EXPECT().LoginWithToken("authcode")
		api.EXPECT().GetSessionKey().Return("mySessionKey")

		secrets := config.NewMockConfig(ctrl)
		gomock.InOrder(
			secrets.EXPECT().Set("session_key", "mySessionKey"),
			secrets.EXPECT().Save(),
		)

		service := &Lastfm{
			api:     api,
			secrets: secrets,
		}

		err := service.Authenticate("authcode")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error for invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := NewMockAPI(ctrl)
		api.EXPECT().LoginWithToken(gomock.Any()).Return(errors.New("error"))

		service := &Lastfm{api: api}

		err := service.Authenticate("authcode")

		if err == nil {
			t.Fatal("Expected an error")
		}
	})

	t.Run("returns error for failed config save", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := NewMockAPI(ctrl)
		api.EXPECT().LoginWithToken(gomock.Any()).Return(nil)
		api.EXPECT().GetSessionKey().Return("mySessionKey")

		secrets := config.NewMockConfig(ctrl)
		secrets.EXPECT().Set(gomock.Any(), gomock.Any())
		secrets.EXPECT().Save().Return(errors.New("save error"))

		service := &Lastfm{
			api:     api,
			secrets: secrets}

		err := service.Authenticate("authcode")

		if err == nil {
			t.Fatal("Expected an error")
		}
	})

	t.Run("returns username from client", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		user := lastfm.UserGetInfo{Name: "Diana"}
		userAPI := NewMockUserAPI(ctrl)
		userAPI.EXPECT().GetInfo(lastfm.P{}).Return(user, nil)

		service := &Lastfm{userAPI: userAPI}

		expected := "Diana"
		got, err := service.GetUsername()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error when failing to read username", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		user := lastfm.UserGetInfo{}
		userAPI := NewMockUserAPI(ctrl)
		userAPI.EXPECT().GetInfo(gomock.Any()).Return(user, errors.New("error"))

		service := &Lastfm{userAPI: userAPI}

		expected := ""
		got, err := service.GetUsername()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if err == nil {
			t.Fatal("Expected an error")
		}
	})

	t.Run("returns map of loved tracks", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		result := lastfm.UserGetLovedTracks{
			Tracks: []struct {
				Name string "xml:\"name\""
				Mbid string "xml:\"mbid\""
				Url  string "xml:\"url\""
				Date struct {
					Uts  string "xml:\"uts,attr\""
					Date string "xml:\",chardata\""
				} "xml:\"date\""
				Artist struct {
					Name string "xml:\"name\""
					Mbid string "xml:\"mbid\""
					Url  string "xml:\"url\""
				} "xml:\"artist\""
				Images []struct {
					Size string "xml:\"size,attr\""
					Url  string "xml:\",chardata\""
				} "xml:\"image\""
				Streamable struct {
					FullTrack  string "xml:\"fulltrack,attr\""
					Streamable string "xml:\",chardata\""
				} "xml:\"streamable\""
			}{
				{
					Name: "Blam (Instrumental)",
					Artist: struct {
						Name string "xml:\"name\""
						Mbid string "xml:\"mbid\""
						Url  string "xml:\"url\""
					}{
						Name: "Awesome Artist",
					},
				},
				{
					Name: "Mr. Testy",
					Artist: struct {
						Name string "xml:\"name\""
						Mbid string "xml:\"mbid\""
						Url  string "xml:\"url\""
					}{
						Name: "Foo & Bar",
					},
				},
			},
		}

		user := lastfm.UserGetInfo{Name: "Diana"}
		userAPI := NewMockUserAPI(ctrl)
		userAPI.EXPECT().GetInfo(lastfm.P{}).Return(user, nil)
		userAPI.EXPECT().GetLovedTracks(lastfm.P{
			"user":  "Diana",
			"limit": 5,
		}).Return(result, nil)

		service := &Lastfm{userAPI: userAPI}

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

	t.Run("returns error when failed to read username for loved tracks", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		user := lastfm.UserGetInfo{}
		userAPI := NewMockUserAPI(ctrl)
		userAPI.EXPECT().GetInfo(lastfm.P{}).Return(user, errors.New("username error"))

		service := &Lastfm{userAPI: userAPI}

		got, err := service.GetLovedTracks(5)

		if err == nil {
			t.Error("Expected an error")
		}

		if len(got) > 0 {
			t.Errorf("Unexpected return value: %v", got)
		}
	})

	t.Run("returns error when failed to read loved tracks", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		user := lastfm.UserGetInfo{Name: "Diana"}
		result := lastfm.UserGetLovedTracks{}
		userAPI := NewMockUserAPI(ctrl)
		userAPI.EXPECT().GetInfo(lastfm.P{}).Return(user, nil)
		userAPI.EXPECT().GetLovedTracks(gomock.Any()).Return(result, errors.New("read error"))

		service := &Lastfm{userAPI: userAPI}

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

		trackAPI := NewMockTrackAPI(ctrl)
		trackAPI.EXPECT().Love(lastfm.P{
			"track":  "Mr. Testy",
			"artist": "Foo & Bar",
		}).Return(nil)

		service := &Lastfm{
			trackAPI: trackAPI,
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

	t.Run("returns error when marking track as loved fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		trackAPI := NewMockTrackAPI(ctrl)
		trackAPI.EXPECT().Love(gomock.Any()).Return(errors.New("api error"))

		service := &Lastfm{
			trackAPI: trackAPI,
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
}

func TestNewLastfm(t *testing.T) {
	t.Run("creates instance when environment is configured", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		secrets := config.NewMockConfig(ctrl)
		secrets.EXPECT().GetString("session_key").Return("mySessionKey")

		os.Setenv("LASTFM_CLIENT_ID", "client_id")
		os.Setenv("LASTFM_CLIENT_SECRET", "client_secret")

		service, err := NewLastfm(secrets)

		if service == nil {
			t.Error("Expected an instance")
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "mySessionKey"
		got := service.api.GetSessionKey()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error when environment is not configured", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		secrets := config.NewMockConfig(ctrl)

		os.Unsetenv("LASTFM_CLIENT_ID")
		os.Unsetenv("LASTFM_CLIENT_SECRET")

		service, err := NewLastfm(secrets)

		if service != nil {
			t.Errorf("Unexpected instance: %v", service)
		}

		if err == nil {
			t.Error("Expected an error")
		}
	})
}
