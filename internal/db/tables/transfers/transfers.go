package transfers

import (
	"database/sql"
	"fmt"
	"log"
)

type sent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type received struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type CoinHistory struct {
	Received *[]received `json:"received"`
	Sent     *[]sent     `json:"sent"`
}

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

func getSenderTransfers(db *sql.DB, employeeId int64) (*[]sent, error) {
	env := "tables.inventory.getSenderTransfers"

	rows, err := db.Query("SELECT to, amount FROM transfers WHERE from = $1;", employeeId)
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}
	defer rows.Close()

	var collection []sent
	for rows.Next() {
		var sentRecord sent
		if err := rows.Scan(&sentRecord.ToUser, &sentRecord.Amount); err != nil {
			log.Printf("%s: failed to get the transfer record, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the transfer record, err: %w", env, err)
		}

		collection = append(collection, sentRecord)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func getReceiverTransfers(db *sql.DB, employeeId int64) (*[]received, error) {
	env := "tables.inventory.getReceiverTransfers"

	rows, err := db.Query("SELECT from, amount FROM transfers WHERE to = $1;", employeeId)
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}
	defer rows.Close()

	var collection []received
	for rows.Next() {
		var receivedRecord received
		if err := rows.Scan(&receivedRecord.FromUser, &receivedRecord.Amount); err != nil {
			log.Printf("%s: failed to get the transfer record, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the transfer record, err: %w", env, err)
		}

		collection = append(collection, receivedRecord)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func GetCoinHistory(db *sql.DB, employeeId int64) (*CoinHistory, error) {
	env := "tables.inventory.GetCoinHistory"

	sent, err := getSenderTransfers(db, employeeId)
	if err != nil {
		log.Printf("%s: failed to get the sent transfers, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to get the sent transfers, err: %w", env, err)
	}

	received, err := getReceiverTransfers(db, employeeId)
	if err != nil {
		log.Printf("%s: failed to get the received transfers, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to get the received transfers, err: %w", env, err)
	}

	return &CoinHistory{Received: received, Sent: sent}, nil
}
