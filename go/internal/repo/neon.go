package repo

import (
	"context"
	"database/sql"
	"fmt"
	"goapp/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type NeonRepo struct {
	db *sql.DB
}

func NewNeonRepo() Repository {
	return &NeonRepo{}
}

func (r *NeonRepo) Connect(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	// probni ping
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}
	r.db = db
	return nil
}

func (r *NeonRepo) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

func (r *NeonRepo) Ping(ctx context.Context) error {
	if r.db == nil {
		return fmt.Errorf("db not connected")
	}
	return r.db.PingContext(ctx)
}

func (r *NeonRepo) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	query := `
		INSERT INTO users.users (email, password, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, email, created_at`
	row := r.db.QueryRowContext(ctx, query, req.Email, req.Password)

	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.CreatedAt); err != nil {
		fmt.Printf("Error scanning row: %v\n", err)
		return nil, err
	}
	return &u, nil
}

func (r *NeonRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, email, password, created_at FROM users.users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *NeonRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password, created_at FROM users.users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *NeonRepo) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) (*models.User, error) {
	query := `UPDATE users.users SET email = $1 WHERE id = $2 RETURNING id, email, created_at`
	row := r.db.QueryRowContext(ctx, query, req.Email, id)
	
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *NeonRepo) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users.users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *NeonRepo) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `SELECT id, email, created_at FROM users.users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *NeonRepo) VerifyPassword(ctx context.Context, email, password string) (*models.User, error) {
	query := `SELECT id, email, created_at FROM users.users WHERE email = $1 AND password = $2`
	row := r.db.QueryRowContext(ctx, query, email, password)
	
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *NeonRepo) UpdatePassword(ctx context.Context, id string, newPassword string) error {
	query := `UPDATE users.users SET password = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newPassword, id)
	return err
}

