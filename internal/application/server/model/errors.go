package model

import (
	"errors"
	"fmt"
)

var (
	UserNotFoundError      = errors.New("user not found")
	UserAlreadyExistsError = errors.New("user is already exists")
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func Error(txt string, code int) error {
	return ErrorResponse{
		Message: txt,
	}
}
