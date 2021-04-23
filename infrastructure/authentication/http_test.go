package authentication

import (
	"net/http/httptest"
	"testing"
)

func TestHttpCallbackProvider(t *testing.T) {
}

func TestHttpCallbackHandler(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	handler := new(httpCallbackHandler)
	handler.ServeHTTP(response, request)
}
