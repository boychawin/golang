package errors

import (
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) error {
	return AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewUnexpectedError(message string) error {
	// fmt.Println(err.Error())
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
