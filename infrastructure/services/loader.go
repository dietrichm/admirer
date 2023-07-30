package services

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/config"
)

type loaderMap map[string]func(secrets config.Config) (domain.Service, error)

type mapServiceLoader struct {
	services     loaderMap
	configLoader config.Loader
}

func (m mapServiceLoader) ForName(serviceName string) (service domain.Service, err error) {
	replaceRegex := regexp.MustCompile("[^a-zA-Z0-9]")
	internalServiceName := strings.ToLower(replaceRegex.ReplaceAllString(serviceName, ""))

	loader, exists := m.services[internalServiceName]

	if !exists {
		return nil, fmt.Errorf("unknown service %q", serviceName)
	}

	secrets, err := m.configLoader.Load("secrets-" + internalServiceName)
	if err != nil {
		return nil, err
	}

	return loader(secrets)
}

func (m mapServiceLoader) Names() (names []string) {
	for name := range m.services {
		names = append(names, name)
	}
	sort.Strings(names)
	return
}
