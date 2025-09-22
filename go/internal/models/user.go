package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` 
	CreatedAt time.Time `json:"created_at"`
}


type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at,omitempty"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	User    *UserResponse `json:"user,omitempty"`
}

type LoginResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	User    *UserResponse `json:"user,omitempty"`
	Token   string        `json:"token,omitempty"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}