package services

import (
	"fmt"

	"github.com/dietrichm/admirer/domain"
)

type loaderMap map[string]func() (domain.Service, error)

// MapServiceLoader loads actual instances of services.
type MapServiceLoader struct {
	services loaderMap
}

// ForName returns service instance for service name.
func (m MapServiceLoader) ForName(serviceName string) (service domain.Service, err error) {
	loader, exists := m.services[serviceName]

	if !exists {
		return nil, fmt.Errorf("unknown service %q", serviceName)
	}

	return loader()
}
