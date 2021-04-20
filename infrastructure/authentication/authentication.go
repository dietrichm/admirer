package authentication

import (
	"context"
	"fmt"
	"net/http"
)

// DefaultCallbackServer is the default authentication callback server.
var DefaultCallbackServer = httpCallbackServer{}

// CallbackServer is the interface for an authentication callback server.
type CallbackServer interface {
	ReadCode(key string) (code string, err error)
}

type httpCallbackServer struct{}

func (h httpCallbackServer) ReadCode(key string) (code string, err error) {
	ctx, closeServer := context.WithCancel(context.Background())

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		code = request.FormValue(key)
		fmt.Fprint(writer, "You may close this window and return to Admirer.")
		closeServer()
	})

	server := &http.Server{
		Addr: ":8080",
	}
	go func() {
		server.ListenAndServe()
	}()

	<-ctx.Done()
	server.Shutdown(ctx)

	if code == "" {
		err = fmt.Errorf("%s missing in request", key)
	}

	return
}
