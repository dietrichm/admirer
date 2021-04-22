package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/authentication"
	"github.com/golang/mock/gomock"
)

func TestLogin(t *testing.T) {
	t.Run("prints service authentication URL and authenticates with received auth code", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		service := domain.NewMockService(ctrl)
		service.EXPECT().Name().AnyTimes().Return("Service")
		service.EXPECT().CreateAuthURL("https://admirer.test").Return("https://service.test/auth")
		service.EXPECT().CodeParam().Return("codeparam")
		service.EXPECT().Authenticate("authcode", "https://admirer.test")
		service.EXPECT().GetUsername().Return("Joe", nil)
		service.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("foobar").Return(service, nil)

		callbackProvider := authentication.NewMockCallbackProvider(ctrl)
		callbackProvider.EXPECT().ReadCode("codeparam").Return("authcode", nil)

		got, err := executeLogin(serviceLoader, callbackProvider, "foobar")
		expected := `Service authentication URL: https://service.test/auth
Logged in on Service as Joe
`

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("authenticates on service with provided auth code", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		service := domain.NewMockService(ctrl)
		service.EXPECT().Authenticate("authcode", "https://admirer.test")
		service.EXPECT().Name().AnyTimes().Return("Service")
		service.EXPECT().GetUsername().Return("Joe", nil)
		service.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		callbackProvider := authentication.NewMockCallbackProvider(ctrl)

		got, err := executeLogin(serviceLoader, callbackProvider, "foobar", "authcode")
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
		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(nil, errors.New(expected))

		callbackProvider := authentication.NewMockCallbackProvider(ctrl)

		output, err := executeLogin(serviceLoader, callbackProvider, "foobar")

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
		service := domain.NewMockService(ctrl)
		service.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(errors.New(expected))
		service.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		callbackProvider := authentication.NewMockCallbackProvider(ctrl)

		output, err := executeLogin(serviceLoader, callbackProvider, "foobar", "authcode")

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
		service := domain.NewMockService(ctrl)
		service.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(nil)
		service.EXPECT().GetUsername().Return("", errors.New(expected))
		service.EXPECT().Close()

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		callbackProvider := authentication.NewMockCallbackProvider(ctrl)

		output, err := executeLogin(serviceLoader, callbackProvider, "foobar", "authcode")

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

func executeLogin(serviceLoader domain.ServiceLoader, callbackProvider authentication.CallbackProvider, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := login(serviceLoader, callbackProvider, buffer, args)
	return buffer.String(), err
}
