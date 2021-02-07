package services

import (
	"fmt"

	"github.com/dietrichm/admirer/services/lastfm"
	"github.com/dietrichm/admirer/services/spotify"
)

// Service is the external service interface.
type Service interface {
	Name() string
	CreateAuthURL() string
	Authenticate(code string) error
	GetUsername() (string, error)
}

// ForName returns service instance for service name.
func ForName(serviceName string) (service Service, err error) {
	switch serviceName {
	case "spotify":
		service, err = spotify.NewSpotify()
	case "lastfm":
		service, err = lastfm.NewLastfm()
	default:
		err = fmt.Errorf("unknown service %q", serviceName)
	}
	return
}
