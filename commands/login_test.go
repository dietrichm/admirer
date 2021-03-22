package commands

import (
	"bytes"
	"testing"

	"github.com/dietrichm/admirer/services"
)

func TestLogin(t *testing.T) {
	t.Run("prints service authentication URL", func(t *testing.T) {
		got, err := executeLogin("foobar")
		expected := "Service authentication URL: https://service.test/auth\n"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func executeLogin(args ...string) (string, error) {
	serviceLoader := new(MockServiceLoader)

	buffer := new(bytes.Buffer)
	loginCommand.SetOutput(buffer)

	err := Login(serviceLoader, loginCommand, args)

	return buffer.String(), err
}

type MockService struct{}

func (m *MockService) Name() string {
	return "Service"
}
func (m *MockService) CreateAuthURL() string {
	return "https://service.test/auth"
}
func (m *MockService) Authenticate(code string) error {
	return nil
}
func (m *MockService) GetUsername() (string, error) {
	return "Joe", nil
}

type MockServiceLoader struct{}

func (m *MockServiceLoader) ForName(serviceName string) (service services.Service, err error) {
	return new(MockService), nil
}
