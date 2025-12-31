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

var (
	ErrInternalServer = errors.New("internal server error")
)

type UseCase interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
}

type useCase struct {
	repo     Repository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewUseCase(repo Repository, log *logrus.Logger, validate *validator.Validate) UseCase {
	return &useCase{repo: repo, log: log, validate: validate}
}

func (u *useCase) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// 1. Validasi Input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("Failed to hash password:", err)
		return nil, ErrInternalServer
	}

	// 3. Construct Entity
	newUser := &User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
		DeletedAt: nil,
	}

	// 4. Simpan ke DB
	if err := u.repo.Save(ctx, newUser); err != nil {
		// Cek error spesifik dari Repository (tanpa variabel ErrUserAlreadyExists lagi)
		if errors.Is(err, ErrEmailTaken) || errors.Is(err, ErrUsernameTaken) {
			return nil, err
		}

		u.log.WithError(err).Error("Failed to save user")
		return nil, ErrInternalServer
	}

	// 5. Return Response
	return &RegisterResponse{
		UserResponse: UserResponse{
			ID:        newUser.ID,
			Username:  newUser.Username,
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt,
		},
	}, nil
}
