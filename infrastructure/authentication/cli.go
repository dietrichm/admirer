package authentication

import (
	"bufio"
	"io"
	"strings"
)

type cliCallbackProvider struct {
	reader io.Reader
}

func (c cliCallbackProvider) ReadCode(key string) (code string, err error) {
	bufferedReader := bufio.NewReader(c.reader)
	code, err = bufferedReader.ReadString('\n')

	if code != "" {
		code = strings.TrimRight(code, "\n")
	}

	return
}
