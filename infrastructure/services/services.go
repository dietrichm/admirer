package services

import (
	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/config"
	"github.com/dietrichm/admirer/infrastructure/services/lastfm"
	"github.com/dietrichm/admirer/infrastructure/services/spotify"
)

// AvailableServices is the configured ServiceLoader for the available services.
var AvailableServices = mapServiceLoader{
	services: loaderMap{
		"spotify": func(secrets config.Config) (domain.Service, error) {
			return spotify.NewSpotify(secrets)
		},
		"lastfm": func(secrets config.Config) (domain.Service, error) {
			return lastfm.NewLastfm(secrets)
		},
	},
	configLoader: config.SecretsLoader,
}
