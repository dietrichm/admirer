package authentication

import (
	"io"
	"net/http"
)

type httpCallbackProvider struct{}

func (h httpCallbackProvider) ReadCode(key string, writer io.Writer) (code string, err error) {
	return
}

type httpCallbackHandler struct{}

func (h httpCallbackHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
}
