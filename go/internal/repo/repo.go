package repo

import (
	"context"
	"goapp/internal/models"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest struktura za ažuriranje korisnika
type UpdateUserRequest struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

type Repository interface {
	Connect(ctx context.Context, dsn string) error
	Close() error
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id int, req UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
	VerifyPassword(ctx context.Context, email, password string) (*models.User, error)
	UpdatePassword(ctx context.Context, id int, newPassword string) error
}