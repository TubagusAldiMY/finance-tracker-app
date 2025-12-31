-- Mengaktifkan ekstensi UUID jika database belum support (opsional di Postgres modern, tapi good practice)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Table: Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY, -- Menggunakan tipe data UUID native
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Index untuk pencarian cepat login berdasarkan email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- 2. Table: Monthly Budgets
-- Penamaan: snake_case (standard PostgreSQL), plural
CREATE TABLE IF NOT EXISTS monthly_budgets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    budget NUMERIC(15, 2) NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
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