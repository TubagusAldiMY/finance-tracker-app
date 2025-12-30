package user

import "time"

type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type RegisterRequest struct {
	ID       string `json:"id"`
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
