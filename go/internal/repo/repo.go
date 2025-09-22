package repo

import (
	"context"
	"goapp/internal/models"
)

type Repository interface {
	Connect(ctx context.Context, dsn string) error
	Close() error
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
	VerifyPassword(ctx context.Context, email, password string) (*models.User, error)
	UpdatePassword(ctx context.Context, id string, newPassword string) error
}