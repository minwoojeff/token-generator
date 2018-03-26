package helpers

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
)

func Open(file string) (*os.File, error) {
	if file == "" {
		return nil, errors.New("Undefined filename")
	}
	dir, _ := os.Getwd()
	f, err := os.Open(dir + file)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Read(reader io.Reader) ([]byte, error) {
	val, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return val, nil
}
