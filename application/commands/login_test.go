package commands

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/config"
	"github.com/golang/mock/gomock"
)

func TestLogin(t *testing.T) {
	t.Run("prints service authentication URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		service := domain.NewMockService(ctrl)
		service.EXPECT().Name().AnyTimes().Return("Service")
		service.EXPECT().CreateAuthURL().Return("https://service.test/auth")

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName("foobar").Return(service, nil)

		secrets := config.NewMockConfig(ctrl)
		configLoader := config.NewMockLoader(ctrl)
		configLoader.EXPECT().Load("secrets").Return(secrets, nil)

		got, err := executeLogin(serviceLoader, configLoader, "foobar")
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

		service := domain.NewMockService(ctrl)
		service.EXPECT().Authenticate("authcode")
		service.EXPECT().AccessToken().Return("access_token")
		service.EXPECT().Name().AnyTimes().Return("Service")
		service.EXPECT().GetUsername().Return("Joe", nil)

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		secrets := config.NewMockConfig(ctrl)
		configLoader := config.NewMockLoader(ctrl)
		configLoader.EXPECT().Load("secrets").Return(secrets, nil)

		gomock.InOrder(
			secrets.EXPECT().Set("service.foobar.access_token", "access_token"),
			secrets.EXPECT().WriteConfig(),
		)

		got, err := executeLogin(serviceLoader, configLoader, "foobar", "authcode")
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

		configLoader := config.NewMockLoader(ctrl)

		output, err := executeLogin(serviceLoader, configLoader, "foobar")

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
		service.EXPECT().Authenticate(gomock.Any()).Return(errors.New(expected))

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		secrets := config.NewMockConfig(ctrl)
		configLoader := config.NewMockLoader(ctrl)
		configLoader.EXPECT().Load(gomock.Any()).Return(secrets, nil)

		output, err := executeLogin(serviceLoader, configLoader, "foobar", "authcode")

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

	t.Run("returns error for failed secrets file write", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		service := domain.NewMockService(ctrl)
		service.EXPECT().Authenticate(gomock.Any()).Return(nil)
		service.EXPECT().AccessToken().AnyTimes().Return("access_token")
		service.EXPECT().GetUsername().AnyTimes().Return("Joe", nil)
		service.EXPECT().Name().AnyTimes().Return("Service")

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		expected := "failed file write"
		secrets := config.NewMockConfig(ctrl)
		secrets.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
		secrets.EXPECT().WriteConfig().Return(errors.New(expected))

		configLoader := config.NewMockLoader(ctrl)
		configLoader.EXPECT().Load(gomock.Any()).Return(secrets, nil)

		output, err := executeLogin(serviceLoader, configLoader, "foobar", "authcode")

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
		service.EXPECT().Authenticate(gomock.Any()).Return(nil)
		service.EXPECT().AccessToken().AnyTimes().Return("access_token")
		service.EXPECT().GetUsername().Return("", errors.New(expected))

		serviceLoader := domain.NewMockServiceLoader(ctrl)
		serviceLoader.EXPECT().ForName(gomock.Any()).Return(service, nil)

		secrets := config.NewMockConfig(ctrl)
		secrets.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
		secrets.EXPECT().WriteConfig().AnyTimes()

		configLoader := config.NewMockLoader(ctrl)
		configLoader.EXPECT().Load(gomock.Any()).Return(secrets, nil)

		output, err := executeLogin(serviceLoader, configLoader, "foobar", "authcode")

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

func executeLogin(serviceLoader domain.ServiceLoader, configLoader config.Loader, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	err := Login(serviceLoader, configLoader, buffer, args)
	return buffer.String(), err
}
