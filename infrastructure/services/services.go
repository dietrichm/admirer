package services

import (
	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/services/lastfm"
	"github.com/dietrichm/admirer/infrastructure/services/spotify"
)

// AvailableServices is the configured ServiceLoader for the available services.
var AvailableServices = MapServiceLoader{
	"spotify": func() (domain.Service, error) {
		return spotify.NewSpotify()
	},
	"lastfm": func() (domain.Service, error) {
		return lastfm.NewLastfm()
	},
}
