package model

import (
	"net/mail"
	"strings"
)

type User struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RegistrationDTO struct {
	User
	RePassword string `json:"re_password,omitempty"`
}

type RegistrationResponse struct {
	User
}

type AuthResponse struct {
	Token string `json:"token"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

func (r *RegistrationDTO) ErrorsRaw(errors map[string]string) string {
	errText := ""

	for _, value := range errors {
		errText += value + "\n"
	}

	return errText
}
