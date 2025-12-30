package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Save(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at, deleted_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.DeletedAt)
	return err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email, password, created_at, deleted_at FROM users WHERE email = $1`

	var user User
	// Scanning
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User tidak ditemukan
		}
		return nil, err
	}

	return &user, nil
}
