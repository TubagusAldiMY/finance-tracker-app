package user

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UseCase interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

type useCase struct {
	repo     Repository
	log      *logrus.Logger
	validate *validator.Validate
	cfg      *viper.Viper
}

func NewUseCase(repo Repository, log *logrus.Logger, validate *validator.Validate, cfg *viper.Viper) UseCase {
	return &useCase{
		repo:     repo,
		log:      log,
		validate: validate,
		cfg:      cfg,
	}
}

// Register Usecase
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

// Login Usecase
func (u *useCase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. Cari User by Email
	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		u.log.WithError(err).Error("Login Failed : Error Finding User")
		return nil, ErrInternalServer
	}

	// 2. Verifikasi Password
	// Jika user nil, return InvalidCredentials
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		u.log.Warnf("Login failed: invalid password for email %s", req.Email)
		return nil, ErrInvalidCredentials
	}

	// 3. Generate JWT Token
	tokenTTL := u.cfg.GetDuration("jwt.ttl")
	if tokenTTL == 0 {
		tokenTTL = 24 * time.Hour // Default value
	}

	// Setup Claims
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"exp":   now.Add(tokenTTL).Unix(),
		"iat":   now.Unix(),
		"name":  user.Username,
		"email": user.Email,
	}

	// PERBAIKAN: Gunakan HS256 untuk secret key string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ambil Secret Key
	jwtSecret := u.cfg.GetString("jwt.secret")
	if jwtSecret == "" {
		u.log.Fatal("JWT Secret Is not Configured")
		return nil, ErrInternalServer
	}

	// Sign Token
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		u.log.WithError(err).Error("Failed To Sign Token")
		return nil, ErrInternalServer
	}

	// 4. Return Response
	return &LoginResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}
