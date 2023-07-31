package authentication

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpCallbackProvider(t *testing.T) {
}

func TestHttpCallbackHandler(t *testing.T) {
	request := httptest.NewRequest("GET", "/?myToken=tokenValue", nil)
	response := httptest.NewRecorder()

	t.Run("saves request form value as specified by injected key", func(t *testing.T) {
		handler := &httpCallbackHandler{
			Key:   "myToken",
			Value: make(chan string, 1),
		}

		go handler.ServeHTTP(response, request)

		assert.Equal(t, "tokenValue", <-handler.Value)
	})

	t.Run("saves empty string when request form value does not exist", func(t *testing.T) {
		handler := &httpCallbackHandler{
			Key:   "nonExisting",
			Value: make(chan string, 1),
		}

		go handler.ServeHTTP(response, request)

		assert.Empty(t, <-handler.Value)
	})
}
