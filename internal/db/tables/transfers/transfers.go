package transfers

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateTransfer(tx *sql.Tx, fromId, toId int64, amount int) error {
	env := "tables.employees.CreateUser"

	stmt, err := tx.Prepare("INSERT INTO transfers(from, to, amount) VALUES($1, $2, $3);")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	_, err = stmt.Exec(fromId, toId, amount)
	if err != nil {
		log.Printf("%s: unmatched arguments to insert, err: %v", env, err)
		return fmt.Errorf("%s: unmatched arguments to insert, err: %w", env, err)
	}

	return nil
}
