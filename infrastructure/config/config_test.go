package config

import (
	"io/ioutil"
	"os"
)

func createFile(contents string) (*os.File, error) {
	file, err := ioutil.TempFile(os.TempDir(), "admirer-")
	if err != nil {
		return nil, err
	}

	if _, err := file.WriteString(contents); err != nil {
		return nil, err
	}

	return file, nil
}
