package services

import (
	"fmt"

	"github.com/dietrichm/admirer/domain"
)

// MapServiceLoader loads actual instances of services.
type MapServiceLoader map[string]func() (domain.Service, error)

// ForName returns service instance for service name.
func (m MapServiceLoader) ForName(serviceName string) (service domain.Service, err error) {
	loader, exists := m[serviceName]

	if !exists {
		return nil, fmt.Errorf("unknown service %q", serviceName)
	}

	return loader()
}
