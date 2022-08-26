package errors

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	Err        error
	StatusCode int
}

func (err *CustomError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err *CustomError) Unwrap() error {
	return err.Err
}

func ParseError(err error) int {
	switch e := err.(type) {
	case *CustomError:
		return e.StatusCode
	default:
		return http.StatusInternalServerError
	}
}
