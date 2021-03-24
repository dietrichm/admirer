package lastfm

import (
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
}
