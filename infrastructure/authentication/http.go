package authentication

import (
	"io"
	"net/http"
)

type httpCallbackProvider struct{}

func (h httpCallbackProvider) ReadCode(key string, writer io.Writer) (code string, err error) {
	handler := &httpCallbackHandler{Key: key}
	server := &http.Server{
		Handler: handler,
	}

	go server.ListenAndServe()

	return
}

type httpCallbackHandler struct {
	Key   string
	Value string
}

func (h *httpCallbackHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.Value = request.FormValue(h.Key)
}
