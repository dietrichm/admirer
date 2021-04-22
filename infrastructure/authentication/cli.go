package authentication

type cliCallbackProvider struct{}

func (c cliCallbackProvider) ReadCode(key string) (code string, err error) {
	return "", nil
}
