package merch

import (
	"database/sql"
	"fmt"
	"log"
)

func GetMerchPrice(db *sql.DB, id int64) (int, error) {
	env := "tables.merch.GetMerchPrice"

	getStmt, err := db.Prepare("SELECT price FROM merch WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the select stmt, err: %v", env, err)
		return 0, fmt.Errorf("%s: failed to prepare the select stmt, err: %w", env, err)
	}

	var price int
	err = getStmt.QueryRow(id).Scan(&price)
	if err != nil {
		log.Printf("%s: failed to get the record, err: %v", env, err)
		return 0, fmt.Errorf("%s: %w", env, err)
	}

	return price, nil
}
