package db

import (
	"database/sql"
	"fmt"
)

func InitDatabase(db *sql.DB) error {
	if err := initMerchTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := insertMerch(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initUsersTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initAccountsTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initTranfersTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initInventoryTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	return nil
}

func initMerchTable(db *sql.DB) error {
	env := "dbinit.initMerchTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS merch(
	    id BIGSERIAL PRIMARY KEY,
	    type VARCHAR(30) NOT NULL,
	    price INTEGER NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func insertMerch(db *sql.DB) error {
	env := "dbinit.insertMerch"

	err := execStatement(db, `
	INSERT INTO merch (type, price)
	VALUES ('t-shirt', 80),
    ('cup', 20),
	('book', 50),
	('pen', 10),
	('powerbank', 200),
	('hoody', 300),
	('umbrella', 200),
	('socks', 10),
	('wallet', 50),
	('pink-hoody', 500)
	ON CONFLICT (type) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initUsersTable(db *sql.DB) error {
	env := "dbinit.initUsersTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS users(
	    id BIGSERIAL PRIMARY KEY,
	    login VARCHAR(30) NOT NULL,
	    password VARCHAR(100) NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	if err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_login ON users(login);"); err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initAccountsTable(db *sql.DB) error {
	env := "dbinit.initAccountsTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS accounts(
	    id BIGSERIAL PRIMARY KEY,
	    user_id BIGINT NOT NULL REFERENCES users(id),
	    coins INTEGER NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	if err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_user ON accounts(user_id);"); err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initTranfersTable(db *sql.DB) error {
	env := "dbinit.initTranfersTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS transfers(
	    id BIGSERIAL PRIMARY KEY,
		from BIGINT NOT NULL REFERENCES users(id)
		to BIGINT NOT NULL REFERENCES users(id)
	    amount INTEGER);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	if err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_sender ON transfers(from);"); err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	if err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_receiver ON transfers(to);"); err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initInventoryTable(db *sql.DB) error {
	env := "dbinit.initInventoryTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS inventory(
	    id BIGSERIAL PRIMARY KEY,
	    employee_id BIGINT NOT NULL REFERENCES users(id),
	    merch_id BIGINT NOT NULL REFERENCES merch(id),
		quantity INTEGER NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	if err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_employee ON inventory(employee_id);"); err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	if err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_merch ON inventory(merch_id);"); err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func execStatement(db *sql.DB, query string) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error occured during preparation: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("error occured during execution: %w", err)
	}

	return nil
}
