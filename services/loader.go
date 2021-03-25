//go:generate mockgen -source ../services/loader.go -destination ../mock_services/loader.go

package services

import (
	"fmt"

	"github.com/dietrichm/admirer/services/lastfm"
	"github.com/dietrichm/admirer/services/spotify"
)

// ServiceLoader loads service instances by name.
type ServiceLoader interface {
	ForName(serviceName string) (Service, error)
}

// MapServiceLoader loads actual instances of services.
type MapServiceLoader struct{}

// ForName returns service instance for service name.
func (d *MapServiceLoader) ForName(serviceName string) (service Service, err error) {
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
