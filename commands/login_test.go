package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dietrichm/admirer/services"
)

func TestLogin(t *testing.T) {
	t.Run("prints service authentication URL", func(t *testing.T) {
		serviceLoader := &MockServiceLoader{new(MockService), nil}

		got, err := executeLogin(serviceLoader, "foobar")
		expected := "Service authentication URL: https://service.test/auth\n"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("authenticates on service with auth code", func(t *testing.T) {
		service := new(MockService)
		serviceLoader := &MockServiceLoader{service, nil}

		got, err := executeLogin(serviceLoader, "foobar", "authcode")
		expected := "Logged in on Service as Joe\n"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if !service.authenticated {
			t.Error("Authenticate() was not called")
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func executeLogin(serviceLoader services.ServiceLoader, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	loginCommand.SetOutput(buffer)

	err := Login(serviceLoader, loginCommand, args)

	return buffer.String(), err
}

type MockService struct {
	authenticated bool
}

func (m *MockService) Name() string {
	return "Service"
}
func (m *MockService) CreateAuthURL() string {
	return "https://service.test/auth"
}
func (m *MockService) Authenticate(code string) error {
	m.authenticated = true
	return nil
}
func (m *MockService) GetUsername() (string, error) {
	return "Joe", nil
}

type MockServiceLoader struct {
	service services.Service
	error   error
}

func (m *MockServiceLoader) ForName(serviceName string) (service services.Service, err error) {
	return m.service, m.error
}
