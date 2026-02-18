package accounts

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/transfers"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/transactions"
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

func transaction(tx *sql.Tx, userId int64, amount int, env string, action func(int, int) int) error {
	getStmt, err := tx.Prepare("SELECT coins FROM accounts WHERE user_id = $1;")
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

	updStmt, err := tx.Prepare("UPDATE merch SET coins = $2 WHERE user_id = $1;")
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

func CreditTo(tx *sql.Tx, userId int64, amount int) error {
	env := "tables.merch.CreditTo"

	return transaction(tx, userId, amount, env, func(a, b int) int {
		return a + b
	})
}

func WriteOff(tx *sql.Tx, userId int64, amount int) error {
	env := "tables.merch.WriteOff"

	return transaction(tx, userId, amount, env, func(a, b int) int {
		return a - b
	})
}

func Transfer(db *sql.DB, userFrom, userTo int64, amount int) error {
	err := transactions.RunInTx(db, func(tx *sql.Tx) error {

		if txErr := WriteOff(tx, userFrom, amount); txErr != nil {
			return txErr
		}

		if txErr := CreditTo(tx, userTo, amount); txErr != nil {
			return txErr
		}

		if txErr := transfers.CreateTransfer(tx, userFrom, userTo, amount); txErr != nil {
			return txErr
		}

		return nil
	})

	return err
}
