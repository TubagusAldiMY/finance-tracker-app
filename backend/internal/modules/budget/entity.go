package budget

import (
	"math"
	"time"
)

type MonthlyBudget struct {
	ID        string    `json:"id" example:"bf8a39e8-4226-4d04-a035-6453181878d6"`
	UserID    string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Budget    float64   `json:"budget" example:"5000000"`
	Date      time.Time `json:"date" example:"2025-01-01T00:00:00Z"`
	CreatedAt time.Time `json:"-"`
}

type CreateBudgetRequest struct {
	Budget float64 `json:"budget" validate:"required,min=1" example:"5000000"`
	Date   string  `json:"date" validate:"required,datetime=2006-01-02" example:"2025-01-01"`
}

type ListBudgetRequest struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

// Helper: Round to 2 decimal places
func RoundFloat(val float64) float64 {
	return math.Round(val*100) / 100
}
