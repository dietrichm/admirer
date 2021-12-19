package authentication

import (
	"net/http/httptest"
	"testing"

	. "github.com/onsi/gomega"
)

func TestHttpCallbackProvider(t *testing.T) {
}

func TestHttpCallbackHandler(t *testing.T) {
	request := httptest.NewRequest("GET", "/?myToken=tokenValue", nil)
	response := httptest.NewRecorder()

	t.Run("saves request form value as specified by injected key", func(t *testing.T) {
		g := NewGomegaWithT(t)

		handler := &httpCallbackHandler{
			Key:   "myToken",
			Value: make(chan string, 1),
		}

		go handler.ServeHTTP(response, request)

		g.Expect(<-handler.Value).To(Equal("tokenValue"))
	})

	t.Run("saves empty string when request form value does not exist", func(t *testing.T) {
		g := NewGomegaWithT(t)

		handler := &httpCallbackHandler{
			Key:   "nonExisting",
			Value: make(chan string, 1),
		}

		go handler.ServeHTTP(response, request)

		g.Expect(<-handler.Value).To(BeEmpty())
	})
}
