package budget

import (
	"math"
	"time"
)

type MonthlyBudget struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Budget    float64   `json:"budget"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"_"`
}

type CreateBudgetRequest struct {
	Budget float64 `json:"budget" validate:"required,min=1"`
	Date   string  `json:"date" validate:"required,datetime=2006-01-02"`
}

type ListBudgetRequest struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

// Helper: Round to 2 decimal places
func RoundFloat(val float64) float64 {
	return math.Round(val*100) / 100
}
