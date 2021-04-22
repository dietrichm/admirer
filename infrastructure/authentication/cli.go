package authentication

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type cliCallbackProvider struct {
	reader io.Reader
}

func (c cliCallbackProvider) ReadCode(key string, writer io.Writer) (code string, err error) {
	fmt.Fprintf(writer, "Please provide %q parameter from the authentication callback URL's query parameters: ", key)

	bufferedReader := bufio.NewReader(c.reader)
	code, err = bufferedReader.ReadString('\n')

	if code != "" {
		code = strings.TrimRight(code, "\n")
	}

	return
}
