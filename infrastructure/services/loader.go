package services

import (
	"fmt"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/config"
)

type loaderMap map[string]func(secrets config.Config) (domain.Service, error)

type mapServiceLoader struct {
	services     loaderMap
	configLoader config.Loader
}

func (m mapServiceLoader) ForName(serviceName string) (service domain.Service, err error) {
	loader, exists := m.services[serviceName]

	if !exists {
		return nil, fmt.Errorf("unknown service %q", serviceName)
	}

	secrets, err := m.configLoader.Load("secrets-" + serviceName)
	if err != nil {
		return nil, err
	}

	return loader(secrets)
}

func (m mapServiceLoader) Names() (names []string) {
	for name := range m.services {
		names = append(names, name)
	}
	return
}
