package budget

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

/*
 * Error Handling
 */
var (
	ErrInternalServer = errors.New("internal server error")
)

type UseCase interface {
	CreateBudget(ctx context.Context, userID string, req *CreateBudgetRequest) (*MonthlyBudget, error)
	GetBudgets(ctx context.Context, userID string, req *ListBudgetRequest) ([]*MonthlyBudget, error)
}

type useCase struct {
	repo     Repository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewUseCase(repo Repository, log *logrus.Logger, validate *validator.Validate) UseCase {
	return &useCase{repo: repo, log: log, validate: validate}
}

func (u useCase) CreateBudget(ctx context.Context, userID string, req *CreateBudgetRequest) (*MonthlyBudget, error) {
	// 1. Validasi Input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2.Parse Date
	parsedDat, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// 3.Buat Entity
	budget := &MonthlyBudget{
		ID:        uuid.New().String(),
		UserID:    userID,
		Budget:    RoundFloat(req.Budget),
		Date:      parsedDat,
		CreatedAt: time.Now(),
	}

	// 4. Simpan Ke DB
	if err := u.repo.Create(ctx, budget); err != nil {
		u.log.WithError(err).Error("Failed to create budget")
		return nil, ErrInternalServer
	}
	return budget, nil
}

func (u useCase) GetBudgets(ctx context.Context, userID string, req *ListBudgetRequest) ([]*MonthlyBudget, error) {
	// Default: Bulan Ini
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Akhir bulan

	// Override jika ada filter
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			startDate = t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			endDate = t
		}
	}

	budgets, err := u.repo.FindMany(ctx, userID, startDate, endDate)
	if err != nil {
		u.log.WithError(err).Error("Failed to fetch budgets")
		return nil, ErrInternalServer
	}

	return budgets, nil
}
