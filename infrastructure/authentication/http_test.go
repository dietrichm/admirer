package authentication

import (
	"net/http/httptest"
	"testing"
)

func TestHttpCallbackProvider(t *testing.T) {
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
}
