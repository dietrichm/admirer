package commands

import (
	"bytes"
	"testing"

	"github.com/dietrichm/admirer/domain"
)

func TestSync(t *testing.T) {
}

func executeSync(serviceLoader domain.ServiceLoader, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := sync(serviceLoader, buffer, args)
	return buffer.String(), err
}
