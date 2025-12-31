package user

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	DeletedAt *time.Time
}

// UserResponse: Format standar data user untuk output JSON
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterRequest: Validasi input saat daftar
type RegisterRequest struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	UserResponse
}

// LoginRequest: Validasi input saat login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse: WAJIB mengandung Token
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	User        UserResponse `json:"user"`
}
