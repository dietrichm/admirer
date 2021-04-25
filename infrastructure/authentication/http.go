package authentication

import (
	"context"
	"io"
	"net/http"
)

type httpCallbackProvider struct {
	server *http.Server
}

func (h *httpCallbackProvider) ReadCode(key string, writer io.Writer) (code string, err error) {
	ctx := context.Background()

	handler := &httpCallbackHandler{Key: key}
	h.server = &http.Server{
		Handler: handler,
	}

	go h.server.ListenAndServe()

	err = h.server.Shutdown(ctx)

	return
}

type httpCallbackHandler struct {
	Key   string
	Value string
}

func (h *httpCallbackHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.Value = request.FormValue(h.Key)
}
