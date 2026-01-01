package user

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailTaken    = errors.New("email already taken")
	ErrUsernameTaken = errors.New("username already taken")
)

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Save(ctx context.Context, user *User) error {
	// Ubah query kolom name -> username
	query := `
		INSERT INTO users (id, username, email, password, created_at, deleted_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.DeletedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		// Cek error Postgres Unique Violation (Code 23505)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Deteksi constraint mana yang kena
			// Pastikan nama constraint di DB Anda sesuai, atau gunakan logic strings.Contains
			if strings.Contains(pgErr.ConstraintName, "email") {
				return ErrEmailTaken
			}
			if strings.Contains(pgErr.ConstraintName, "username") {
				return ErrUsernameTaken
			}
		}
		return err
	}
	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password, created_at, deleted_at FROM users WHERE email = $1`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, username, email, password, created_at, deleted_at FROM users WHERE id = $1`

	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
