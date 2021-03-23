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

		authcode := "authcode"
		got, err := executeLogin(serviceLoader, "foobar", authcode)
		expected := "Logged in on Service as Joe\n"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		got = service.authenticatedWith
		if got != authcode {
			t.Errorf("expected %q, got %q", authcode, got)
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error for unknown service", func(t *testing.T) {
		expected := "unknown service"
		serviceLoader := &MockServiceLoader{nil, errors.New(expected)}

		output, err := executeLogin(serviceLoader, "foobar")

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}

		if err == nil {
			t.Error("Expected an error")
		}

		got := err.Error()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error for failed authentication", func(t *testing.T) {
		expected := "failed authentication"
		service := &MockService{errors.New(expected), ""}
		serviceLoader := &MockServiceLoader{service, nil}

		output, err := executeLogin(serviceLoader, "foobar", "authcode")

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}

		if err == nil {
			t.Error("Expected an error")
		}

		got := err.Error()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
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
	authenticationError error
	usernameError       error
	authenticatedWith   string
}

func (m *MockService) Name() string {
	return "Service"
}
func (m *MockService) CreateAuthURL() string {
	return "https://service.test/auth"
}
func (m *MockService) Authenticate(code string) error {
	m.authenticatedWith = code
	return m.authenticationError
}
func (m *MockService) GetUsername() (string, error) {
	return "Joe", m.usernameError
}

type MockServiceLoader struct {
	service services.Service
	error   error
}

func (m *MockServiceLoader) ForName(serviceName string) (service services.Service, err error) {
	return m.service, m.error
}
