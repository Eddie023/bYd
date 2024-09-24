package web

import "errors"

type shutdownError struct {
	Message string
}

// Error is the implementation of the error interface.
func (se *shutdownError) Error() string {
	return se.Message
}

func NewshutdownError(msg string) error {
	return &shutdownError{
		Message: msg,
	}
}

func IsShutdown(err error) bool {
	var se *shutdownError

	return errors.As(err, &se)
}
