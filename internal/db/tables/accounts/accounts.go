package accounts

import (
	"database/sql"
	"fmt"
	"log"
)

func GetBalanceById(db *sql.DB, userId int64) (int, error) {
	env := "tables.merch.GetBalanceById"

	getStmt, err := db.Prepare("SELECT coins FROM accounts WHERE user_id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the select stmt, err: %v", env, err)
		return 0, fmt.Errorf("%s: failed to prepare the select stmt, err: %w", env, err)
	}

	var balance int
	err = getStmt.QueryRow(userId).Scan(&balance)
	if err != nil {
		log.Printf("%s: failed to get the record, err: %v", env, err)
		return 0, fmt.Errorf("%s: %w", env, err)
	}

	return balance, nil
}

func transaction(db *sql.DB, userId int64, amount int, env string, action func(int, int) int) error {
	getStmt, err := db.Prepare("SELECT coins FROM accounts WHERE user_id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the select stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the select stmt, err: %w", env, err)
	}

	var balance int
	err = getStmt.QueryRow(userId).Scan(&balance)
	if err != nil {
		log.Printf("%s: failed to get the record, err: %v", env, err)
		return fmt.Errorf("%s: %w", env, err)
	}

	balance = action(balance, amount)

	updStmt, err := db.Prepare("UPDATE merch SET coins = $2 WHERE user_id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the update stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the update stmt, err: %w", env, err)
	}

	_, err = updStmt.Exec(userId, balance)
	if err != nil {
		log.Printf("%s: failed to execute the update stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to execute the update stmt, err: %w", env, err)
	}

	return nil
}

func CreditTo(db *sql.DB, userId int64, amount int) error {
	env := "tables.merch.CreditTo"

	return transaction(db, userId, amount, env, func(a, b int) int {
		return a + b
	})
}

func WriteOff(db *sql.DB, userId int64, amount int) error {
	env := "tables.merch.WriteOff"

	return transaction(db, userId, amount, env, func(a, b int) int {
		return a - b
	})
}
