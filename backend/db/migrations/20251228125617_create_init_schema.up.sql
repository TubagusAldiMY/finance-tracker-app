-- Mengaktifkan ekstensi UUID jika database belum support (opsional di Postgres modern, tapi good practice)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Table: Users
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY, -- Menggunakan tipe data UUID native
                                     name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Nullable untuk soft delete
                             );

-- Index untuk pencarian cepat login berdasarkan email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 2. Table: Monthly Budgets
-- Penamaan: snake_case (standard PostgreSQL), plural
CREATE TABLE IF NOT EXISTS monthly_budgets (
                                               id UUID PRIMARY KEY,
                                               user_id UUID NOT NULL,
                                               budget NUMERIC(15, 2) NOT NULL, -- Numeric untuk uang (15 digit, 2 desimal)
    date TIMESTAMP WITH TIME ZONE NOT NULL, -- Menyimpan periode budget (biasanya awal bulan)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

                             CONSTRAINT fk_user
                             FOREIGN KEY(user_id)
    REFERENCES users(id)
                         ON DELETE CASCADE -- Hapus user = hapus budget
    );

-- Index untuk mengambil list budget berdasarkan user dan rentang tanggal
CREATE INDEX IF NOT EXISTS idx_budgets_user_date ON monthly_budgets(user_id, date);

-- 3. Table: Histories (Pengeluaran)
CREATE TABLE IF NOT EXISTS histories (
                                         id UUID PRIMARY KEY,
                                         budget_id UUID NOT NULL,
                                         date TIMESTAMP WITH TIME ZONE NOT NULL,
                                         amount NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

                             CONSTRAINT fk_budget
                             FOREIGN KEY(budget_id)
    REFERENCES monthly_budgets(id)
                         ON DELETE CASCADE -- Hapus budget = hapus history
    );

-- Index untuk load history per budget
CREATE INDEX IF NOT EXISTS idx_histories_budget ON histories(budget_id);