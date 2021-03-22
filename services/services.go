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

// ServiceLoader loads service instances by name.
type ServiceLoader interface {
	ForName(serviceName string) (Service, error)
}

// DefaultServiceLoader loads actual instances of services.
type DefaultServiceLoader struct{}

// ForName returns service instance for service name.
func (d *DefaultServiceLoader) ForName(serviceName string) (service Service, err error) {
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
