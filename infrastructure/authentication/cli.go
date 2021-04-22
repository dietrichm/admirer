package authentication

import "io"

type cliCallbackProvider struct {
	reader io.Reader
}

func (c cliCallbackProvider) ReadCode(key string) (code string, err error) {
	return "", nil
}
