package model

import (
	"errors"
	"fmt"
)

var (
	// UserNotFoundError возникает, если пользователь не найден в базе данных.
	UserNotFoundError = errors.New("user not found")
	// UserAlreadyExistsError возникает при попытке регистрации пользователя с уже существующим email.
	UserAlreadyExistsError = errors.New("user is already exists")
	// DataNotChangedError указывает на отсутствие изменений в редактируемых данных.
	DataNotChangedError = errors.New("data not changed")
	// AccessDeniedError возвращается при попытке доступа к чужим данным.
	AccessDeniedError = errors.New("access denied")
	// AuthorizationError указывает на ошибку авторизации или генерации токена.
	AuthorizationError = errors.New("authorization error")
	// PasswordNotMatchError возникает, если пароли и их подтверждение не совпадают.
	PasswordNotMatchError = errors.New("password not match")
	// LoginAndPasswordError указывает на некорректную пару логин-пароль.
	LoginAndPasswordError = errors.New("not correct login or password")
	// NotFoundError общий тип ошибки для отсутствующих данных.
	NotFoundError = errors.New("not found")
)

// ErrorResponse представляет стандартную структуру ошибки для ответа API.
type ErrorResponse struct {
	Message string `json:"message"` // Описание ошибки
}

// ValidationErrorResponse используется для отправки ошибок валидации формы.
//
// Расширяет ErrorResponse и добавляет карту ошибок по полям.
type ValidationErrorResponse struct {
	ErrorResponse
	Errors map[string]string `json:"errors"` // Карта ошибок валидации по именам полей
}

// Error реализует интерфейс error для ErrorResponse.
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

// Error создаёт новый ErrorResponse по заданному сообщению и коду ошибки.
//
// Аргумент `code` не используется в структуре, но может быть полезен для дальнейшего расширения.
func Error(txt string, code int) error {
	return ErrorResponse{
		Message: txt,
	}
}
