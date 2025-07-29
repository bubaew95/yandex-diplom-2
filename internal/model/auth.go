package model

import (
	"net/mail"
	"strings"
)

// User представляет зарегистрированного пользователя.
//
// Используется в базе данных, а также в обмене данными между слоями.
type User struct {
	ID        int64  `json:"id,omitempty"` // Уникальный идентификатор пользователя
	FirstName string `json:"first_name"`   // Имя
	LastName  string `json:"last_name"`    // Фамилия
	Email     string `json:"email"`        // Email (используется как логин)
	Password  string `json:"password"`     // Хэш пароля
}

// RegistrationDTO содержит данные, необходимые для регистрации пользователя,
// включая подтверждение пароля.
type RegistrationDTO struct {
	User
	RePassword string `json:"re_password,omitempty"` // Подтверждение пароля
}

// RegistrationResponse представляет ответ на регистрацию пользователя.
type RegistrationResponse struct {
	User
}

// AuthResponse содержит токен, выданный после успешной авторизации или регистрации.
type AuthResponse struct {
	Token string `json:"token"`
}

// LoginDTO содержит данные, необходимые для входа пользователя.
type LoginDTO struct {
	Email    string `json:"email"`    // Email
	Password string `json:"password"` // Пароль в открытом виде (будет проверен и захеширован)
}

// Validate выполняет базовую валидацию полей RegistrationDTO.
//
// Проверяется:
//   - наличие email и его корректность,
//   - наличие имени и фамилии,
//   - наличие пароля и совпадение с подтверждением.
//
// Возвращает map[string]string с названиями полей и текстами ошибок.
func (r *RegistrationDTO) Validate() map[string]string {
	fields := make(map[string]string)

	if r.Email == "" {
		fields["email"] = "email is required"
	} else {
		_, err := mail.ParseAddress(r.Email)
		if err != nil {
			fields["email"] = err.Error()
		}
	}

	if r.Password == "" {
		fields["password"] = "password is required"
	}

	if r.Password != r.RePassword {
		fields["password"] = "passwords not match"
	}

	if strings.TrimSpace(r.LastName) == "" {
		fields["last_name"] = "last name is required"
	}

	if strings.TrimSpace(r.FirstName) == "" {
		fields["first_name"] = "first name is required"
	}

	return fields
}

// ErrorsRaw преобразует карту ошибок в RegistrationDTO в многострочную строку,
// где каждая ошибка на новой строке.
func (r *RegistrationDTO) ErrorsRaw(errors map[string]string) string {
	errText := ""

	for _, value := range errors {
		errText += value + "\n"
	}

	return errText
}
