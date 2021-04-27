package authentication

import (
	"net/http/httptest"
	"testing"
)

func TestHttpCallbackProvider(t *testing.T) {
	provider := new(httpCallbackProvider)
	_, err := provider.ReadCode("myToken", nil)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestHttpCallbackHandler(t *testing.T) {
	request := httptest.NewRequest("GET", "/?myToken=tokenValue", nil)
	response := httptest.NewRecorder()

	t.Run("saves request form value as specified by injected key and call done function", func(t *testing.T) {
		called := false
		doneFunc := func() { called = true }

		handler := httpCallbackHandler{Key: "myToken", DoneFunc: doneFunc}
		handler.ServeHTTP(response, request)

		got := handler.Value
		expected := "tokenValue"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if !called {
			t.Error("Done func was not called")
		}
	})

	t.Run("saves empty string when request form value does not exist", func(t *testing.T) {
		called := false
		doneFunc := func() { called = true }

		handler := httpCallbackHandler{Key: "nonExisting", DoneFunc: doneFunc}
		handler.ServeHTTP(response, request)

		got := handler.Value
		expected := ""

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}

		if !called {
			t.Error("Done func was not called")
		}
	})
}
