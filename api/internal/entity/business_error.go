package entity

import "fmt"

type BusinessError struct {
	err error
}

func NewBusinessError(format string, a ...any) *BusinessError {
	return &BusinessError{err: fmt.Errorf(format, a...)}
}

func (e *BusinessError) Error() string {
	return e.err.Error()
}
