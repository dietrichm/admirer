package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dietrichm/admirer/mock_services"
	"github.com/dietrichm/admirer/services"
	"github.com/golang/mock/gomock"
)

func TestLogin(t *testing.T) {
	t.Run("prints service authentication URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		service := mock_services.NewMockService(ctrl)
		service.EXPECT().Name().AnyTimes().Return("Service")
		service.EXPECT().CreateAuthURL().Return("https://service.test/auth")

		serviceLoader := mock_services.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("foobar").Return(service, nil)

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
		ctrl := gomock.NewController(t)

		service := mock_services.NewMockService(ctrl)
		service.EXPECT().Authenticate("authcode")
		service.EXPECT().Name().AnyTimes().Return("Service")
		service.EXPECT().GetUsername().Return("Joe", nil)

		serviceLoader := mock_services.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		got, err := executeLogin(serviceLoader, "foobar", "authcode")
		expected := "Logged in on Service as Joe\n"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("returns error for unknown service", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		expected := "unknown service"
		serviceLoader := mock_services.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(nil, errors.New(expected))

		output, err := executeLogin(serviceLoader, "foobar")

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}

		if err == nil {
			t.Fatal("Expected an error")
		}

		got := err.Error()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error for failed authentication", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		expected := "failed authentication"
		service := mock_services.NewMockService(ctrl)
		service.EXPECT().Authenticate(gomock.Any()).Return(errors.New(expected))

		serviceLoader := mock_services.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		output, err := executeLogin(serviceLoader, "foobar", "authcode")

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}

		if err == nil {
			t.Fatal("Expected an error")
		}

		got := err.Error()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("returns error for failed username retrieval", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		expected := "failed username retrieval"
		service := mock_services.NewMockService(ctrl)
		service.EXPECT().Authenticate(gomock.Any()).Return(nil)
		service.EXPECT().GetUsername().Return("", errors.New(expected))

		serviceLoader := mock_services.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		output, err := executeLogin(serviceLoader, "foobar", "authcode")

		if output != "" {
			t.Errorf("Unexpected output: %v", output)
		}

		if err == nil {
			t.Fatal("Expected an error")
		}

		got := err.Error()

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

func executeLogin(serviceLoader services.ServiceLoader, args ...string) (string, error) {
	buffer := new(bytes.Buffer)

	err := Login(serviceLoader, buffer, args)

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
