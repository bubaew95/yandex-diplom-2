package model

type User struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RegistrationRequest struct {
	User
	RePassword string `json:"re_password,omitempty"`
}

type RegistrationResponse struct {
	User
}

type AuthResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
