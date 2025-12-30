package user

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UseCase interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
}

type useCase struct {
	repo     Repository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewUseCase(repo Repository, log *logrus.Logger, validate *validator.Validate) UseCase {
	return &useCase{repo: repo, log: log, validate: validate}
}

func (u *useCase) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	// 1. Validasi Input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Cek Email Duplikat
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exist, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		u.log.Error("Database error checking email:", err)
		return nil, errors.New("internal server error")
	}
	if exist != nil {
		return nil, errors.New("email already exists")
	}

	// 3. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("Failed to hash password:", err)
		return nil, errors.New("failed to process request")
	}
	// 4. Construct Entity
	newUser := &User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
		DeletedAt: nil,
	}

	// 5. Simpan ke DB
	if err := u.repo.Save(ctx, newUser); err != nil {
		u.log.Error("Failed to save user:", err)
		return nil, errors.New("failed to register user")
	}

	return newUser, nil
}
