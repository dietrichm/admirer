//go:generate mockgen -source services.go -destination services_mock.go -package services

package services

import (
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

// AvailableServices is the configured ServiceLoader for the available services.
var AvailableServices = MapServiceLoader{
	"spotify": func() (Service, error) {
		return spotify.NewSpotify()
	},
	"lastfm": func() (Service, error) {
		return lastfm.NewLastfm()
	},
}
