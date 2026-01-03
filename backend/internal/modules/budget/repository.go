package budget

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, budget *MonthlyBudget) error
	FindMany(ctx context.Context, userID string, startDate, endDate time.Time) ([]*MonthlyBudget, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

/*
* CREATE data anggaran bulanan baru (MonthlyBudget) ke dalam tabel monthly_budgets di database
* @param ctx context.Context = Context untuk mengatur batas waktu (timeout) atau pembatalan request.
* @param b *MonthlyBudget = Pointer ke struct MonthlyBudget yang berisi data yang akan disimpan.
* @return error = Mengembalikan nil jika berhasil, atau pesan kesalahan jika operasi database gagal.
 */
func (r repository) Create(ctx context.Context, b *MonthlyBudget) error {
	query := `
		INSERT INTO monthly_budgets (id, user_id, budget, date, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		`

	_, err := r.db.Exec(ctx, query, b.ID, b.UserID, b.Budget, b.Date, b.CreatedAt)
	return err
}

/*
 * Mencari data anggaran bulanan berdasarkan ID user dan tanggal awal dan akhir.
 * @param ctx context.Context = Context untuk mengatur batas waktu (timeout) atau pembatalan request.
 * @param userID string = ID user yang ingin mendapatkan data anggaran bulanan.
 * @param startDate, endDate time.Time = Tanggal awal dan akhir untuk mencari data anggaran bulanan.
 * @return []*MonthlyBudget = Array dari struct MonthlyBudget yang berisi data anggaran bulanan yang ditemukan.
 */
func (r *repository) FindMany(ctx context.Context, userID string, startDate, endDate time.Time) ([]*MonthlyBudget, error) {
	query := `
		SELECT id, user_id, budget, date, created_at
		FROM monthly_budgets
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []*MonthlyBudget
	for rows.Next() {
		var b MonthlyBudget
		if err := rows.Scan(&b.ID, &b.UserID, &b.Budget, &b.Date, &b.CreatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, &b)
	}

	return budgets, nil
}
