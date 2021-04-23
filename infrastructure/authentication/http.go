package authentication

import "io"

type httpCallbackProvider struct{}

func (h httpCallbackProvider) ReadCode(key string, writer io.Writer) (code string, err error) {
	return
}
