//go:generate mockgen -source ../services/loader.go -destination ../mock_services/loader.go

package services

import "fmt"

// ServiceLoader loads service instances by name.
type ServiceLoader interface {
	ForName(serviceName string) (Service, error)
}

// MapServiceLoader loads actual instances of services.
type MapServiceLoader map[string]func() (Service, error)

// ForName returns service instance for service name.
func (m MapServiceLoader) ForName(serviceName string) (service Service, err error) {
	loader, exists := m[serviceName]

	if !exists {
		return nil, fmt.Errorf("unknown service %q", serviceName)
	}

	return loader()
}
