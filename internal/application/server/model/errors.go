package model

import (
	"errors"
	"fmt"
)

var (
	UserNotFoundError      = errors.New("user not found")
	UserAlreadyExistsError = errors.New("user is already exists")
	DataNotChangedError    = errors.New("data not changed")
	AccessDeniedError      = errors.New("access denied")
	AuthorizationError     = errors.New("authorization error")
	PasswordNotMatchError  = errors.New("password not match")
	LoginAndPasswordError  = errors.New("not correct login or password")
	NotFoundError          = errors.New("not found")
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	ErrorResponse
	Errors map[string]string `json:"errors"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func Error(txt string, code int) error {
	return ErrorResponse{
		Message: txt,
	}
}

func ValidationError(errors interface{}) error {
	return ValidationErrorResponse{
		ErrorResponse: ErrorResponse{
			Message: "validation error",
		},
		Errors: errors.(map[string]string),
	}
}
