package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/dietrichm/admirer/services"
)

func TestReturnsServiceAuthenticationUrl(t *testing.T) {
	os.Setenv("SPOTIFY_CLIENT_ID", "foo")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "bar")

	got, err := executeLogin("spotify")
	expected := "Spotify authentication URL: https://"

	if !strings.Contains(got, expected) {
		t.Errorf("expected %q, got %q", expected, got)
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func executeLogin(args ...string) (string, error) {
	serviceLoader := new(services.DefaultServiceLoader)

	buffer := new(bytes.Buffer)
	loginCommand.SetOutput(buffer)

	err := Login(serviceLoader, loginCommand, args)

	return buffer.String(), err
}
