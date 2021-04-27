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
	ctx, cancel := context.WithCancel(context.Background())

	handler := &httpCallbackHandler{Key: key, DoneFunc: cancel}
	h.server = &http.Server{
		Handler: handler,
	}

	go h.server.ListenAndServe()

	err = h.server.Shutdown(ctx)

	return
}

type httpCallbackHandler struct {
	Key      string
	Value    string
	DoneFunc func()
}

func (h *httpCallbackHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.Value = request.FormValue(h.Key)
	h.DoneFunc()
}
