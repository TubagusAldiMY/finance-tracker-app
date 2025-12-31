package user_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ==========================================
// 1. MOCK OBJECTS
// ==========================================

// MockRepository memalsukan behavior Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	// Handle jika return nil (User not found)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// ==========================================
// 2. HELPER SETUP
// ==========================================

// setupTest mengembalikan Interface UseCase, MockRepo, dan Config untuk dimanipulasi
func setupTest() (user.UseCase, *MockRepository, *viper.Viper) {
	mockRepo := new(MockRepository)

	// Logger buang ke tong sampah (supaya terminal bersih)
	log := logrus.New()
	log.SetOutput(io.Discard)

	validate := validator.New()

	// Config default yang valid
	cfg := viper.New()
	cfg.Set("jwt.secret", "secret_key_testing_123")
	cfg.Set("jwt.ttl", "1h")

	useCase := user.NewUseCase(mockRepo, log, validate, cfg)

	return useCase, mockRepo, cfg
}

// ==========================================
// 3. GROUP: REGISTER TESTS
// ==========================================

func TestRegister_Success(t *testing.T) {
	u, mockRepo, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "validuser",
		Email:    "valid@example.com",
		Password: "password123",
	}

	// Expectation: Repo.Save dipanggil sekali dengan data yang cocok
	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(userObj *user.User) bool {
		return userObj.Email == req.Email && userObj.Username == req.Username && userObj.ID != ""
	})).Return(nil)

	// Action
	resp, err := u.Register(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.Email)
	assert.NotEmpty(t, resp.ID)

	mockRepo.AssertExpectations(t)
}

func TestRegister_ValidationError(t *testing.T) {
	u, mockRepo, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "",                  // Invalid: Kosong
		Email:    "invalid-email-fmt", // Invalid: Format salah
		Password: "123",               // Invalid: Terlalu pendek
	}

	// Expectation: Repo.Save TIDAK BOLEH dipanggil
	// (Karena harusnya gagal di layer validasi validator/struct)

	// Action
	resp, err := u.Register(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertNotCalled(t, "Save") // Pastikan tidak tembus ke DB
}

func TestRegister_DuplicateEmail(t *testing.T) {
	u, mockRepo, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "newuser",
		Email:    "taken@example.com",
		Password: "password123",
	}

	// Expectation: Repo return error EmailTaken
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(user.ErrEmailTaken)

	resp, err := u.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrEmailTaken, err)
	assert.Nil(t, resp)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	u, mockRepo, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "takenuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	// Expectation: Repo return error UsernameTaken
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(user.ErrUsernameTaken)

	resp, err := u.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrUsernameTaken, err)
	assert.Nil(t, resp)
}

func TestRegister_RepositoryError(t *testing.T) {
	u, mockRepo, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "user",
		Email:    "email@example.com",
		Password: "password123",
	}

	// Expectation: Repo gagal koneksi DB
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db connection lost"))

	resp, err := u.Register(context.Background(), req)

	// Assertions: Harus return ErrInternalServer (jangan bocorkan error asli ke user)
	assert.Error(t, err)
	assert.Equal(t, user.ErrInternalServer, err)
	assert.Nil(t, resp)
}

// ==========================================
// 4. GROUP: LOGIN TESTS
// ==========================================

func TestLogin_Success(t *testing.T) {
	u, mockRepo, _ := setupTest()

	// Siapkan password hash yang valid
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	dummyUser := &user.User{
		ID:        "uuid-123",
		Username:  "loginuser",
		Email:     "login@example.com",
		Password:  string(hashedPwd),
		CreatedAt: time.Now(),
	}

	req := &user.LoginRequest{
		Email:    "login@example.com",
		Password: "password123", // Password raw yang cocok
	}

	// Expectation: FindByEmail sukses
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)

	// Action
	resp, err := u.Login(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken) // Token harus ada
	assert.Equal(t, "Bearer", resp.TokenType)
	assert.Equal(t, dummyUser.Email, resp.User.Email)
}

func TestLogin_UserNotFound(t *testing.T) {
	u, mockRepo, _ := setupTest()

	req := &user.LoginRequest{
		Email:    "ghost@example.com",
		Password: "password123",
	}

	// Expectation: FindByEmail return nil user & nil error (atau error not found, tergantung implementasi repo)
	// Sesuai implementasi mock kita: return nil, nil
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)

	resp, err := u.Login(context.Background(), req)

	// Harus InvalidCredentials, JANGAN "User Not Found" (Security)
	assert.Error(t, err)
	assert.Equal(t, user.ErrInvalidCredentials, err)
	assert.Nil(t, resp)
}

func TestLogin_WrongPassword(t *testing.T) {
	u, mockRepo, _ := setupTest()

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("realpassword"), bcrypt.DefaultCost)
	dummyUser := &user.User{
		Email:    "user@example.com",
		Password: string(hashedPwd),
	}

	req := &user.LoginRequest{
		Email:    "user@example.com",
		Password: "WRONG_PASSWORD",
	}

	// Expectation: User ketemu, tapi nanti bcrypt compare gagal
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)

	resp, err := u.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrInvalidCredentials, err)
	assert.Nil(t, resp)
}

func TestLogin_RepositoryError(t *testing.T) {
	u, mockRepo, _ := setupTest()

	req := &user.LoginRequest{
		Email:    "error@example.com",
		Password: "any",
	}

	// Expectation: DB mati
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, errors.New("db timeout"))

	resp, err := u.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrInternalServer, err)
	assert.Nil(t, resp)
}

func TestLogin_MissingJWTSecret(t *testing.T) {
	// Setup Manual Khusus case ini (karena butuh config rusak)
	mockRepo := new(MockRepository)
	log := logrus.New()
	log.SetOutput(io.Discard)
	validate := validator.New()

	// CONFIG RUSAK (Secret Kosong)
	cfg := viper.New()
	cfg.Set("jwt.secret", "")

	u := user.NewUseCase(mockRepo, log, validate, cfg)

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	dummyUser := &user.User{
		Email:    "test@example.com",
		Password: string(hashedPwd),
	}

	req := &user.LoginRequest{Email: "test@example.com", Password: "pass"}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)

	// Action
	resp, err := u.Login(context.Background(), req)

	// Expectation: Harusnya error internal server (daripada panic)
	assert.Error(t, err)
	assert.Equal(t, user.ErrInternalServer, err)
	assert.Nil(t, resp)
}
