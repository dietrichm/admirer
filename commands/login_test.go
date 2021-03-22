package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestReturnsServiceAuthenticationUrl(t *testing.T) {
	os.Setenv("SPOTIFY_CLIENT_ID", "foo")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "bar")

	rootCommand.SetArgs([]string{"login", "spotify"})

	got, err := executeCommand(rootCommand)
	expected := "Spotify authentication URL: https://"

	if !strings.Contains(got, expected) {
		t.Errorf("expected %q, got %q", expected, got)
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func executeCommand(command *cobra.Command) (string, error) {
	buffer := new(bytes.Buffer)
	command.SetOutput(buffer)

	_, err := command.ExecuteC()

	return buffer.String(), err
}
