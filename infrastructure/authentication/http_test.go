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

	t.Run("saves request form value as specified by injected key and call done function", func(t *testing.T) {
		g := NewGomegaWithT(t)

		called := false
		doneFunc := func() { called = true }

		handler := httpCallbackHandler{Key: "myToken", DoneFunc: doneFunc}
		handler.ServeHTTP(response, request)

		g.Expect(handler.Value).To(Equal("tokenValue"))
		g.Expect(called).To(BeTrue())
	})

	t.Run("saves empty string when request form value does not exist", func(t *testing.T) {
		g := NewGomegaWithT(t)

		called := false
		doneFunc := func() { called = true }

		handler := httpCallbackHandler{Key: "nonExisting", DoneFunc: doneFunc}
		handler.ServeHTTP(response, request)

		g.Expect(handler.Value).To(BeEmpty())
		g.Expect(called).To(BeTrue())
	})
}
