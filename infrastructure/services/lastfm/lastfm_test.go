package lastfm

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shkh/lastfm-go/lastfm"
)

func TestLastfm(t *testing.T) {
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

	t.Run("authenticates using authorization code", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		api := NewMockAPI(ctrl)
		api.EXPECT().LoginWithToken("authcode")

		service := &Lastfm{api: api}

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
}

func TestNewLastfm(t *testing.T) {
	t.Run("creates instance when environment is configured", func(t *testing.T) {
		os.Setenv("LASTFM_CLIENT_ID", "client_id")
		os.Setenv("LASTFM_CLIENT_SECRET", "client_secret")

		service, err := NewLastfm("")

		if service == nil {
			t.Error("Expected an instance")
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error when environment is not configured", func(t *testing.T) {
		os.Unsetenv("LASTFM_CLIENT_ID")
		os.Unsetenv("LASTFM_CLIENT_SECRET")

		service, err := NewLastfm("")

		if service != nil {
			t.Errorf("Unexpected instance: %v", service)
		}

		if err == nil {
			t.Error("Expected an error")
		}
	})
}
