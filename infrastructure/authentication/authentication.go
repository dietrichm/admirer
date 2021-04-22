//go:generate mockgen -source authentication.go -destination authentication_mock.go -package authentication

package authentication

import (
	"io"
	"os"
)

// DefaultCallbackProvider is the default callback provider for services.
var DefaultCallbackProvider = &cliCallbackProvider{os.Stdin}

// CallbackProvider provides a callback mechanism for authenticating services.
type CallbackProvider interface {
	ReadCode(key string, writer io.Writer) (code string, err error)
}
