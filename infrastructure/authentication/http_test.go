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

	t.Run("saves request form value as specified by injected key", func(t *testing.T) {
		handler := httpCallbackHandler{Key: "myToken"}
		handler.ServeHTTP(response, request)

		got := handler.Value
		expected := "tokenValue"

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("saves empty string when request form value does not exist", func(t *testing.T) {
		handler := httpCallbackHandler{Key: "nonExisting"}
		handler.ServeHTTP(response, request)

		got := handler.Value
		expected := ""

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}
