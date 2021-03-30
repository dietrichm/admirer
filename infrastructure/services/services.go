package services

import (
	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/config"
	"github.com/dietrichm/admirer/infrastructure/services/lastfm"
	"github.com/dietrichm/admirer/infrastructure/services/spotify"
)

// AvailableServices is the configured ServiceLoader for the available services.
var AvailableServices = MapServiceLoader{
	services: loaderMap{
		"spotify": func() (domain.Service, error) {
			return spotify.NewSpotify()
		},
		"lastfm": func() (domain.Service, error) {
			secrets, err := config.ConfigLoader.Load("secrets")
			if err != nil {
				return nil, err
			}

			return lastfm.NewLastfm(secrets.GetString("service.lastfm.access_token"))
		},
	},
	configLoader: config.ConfigLoader,
}
