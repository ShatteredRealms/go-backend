package mocks

import "errors"

type ErrReader int

func (ErrReader) Read(p []byte) (int, error) {
	return 0, errors.New("test error")
}
