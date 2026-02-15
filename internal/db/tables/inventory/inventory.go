package inventory

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Inventory struct {
	Id         int64 `json:"id"`
	EmployeeId int64 `json:"employee_id"`
	MerchId    int64 `json:"merch_id"`
	Quantity   int   `json:"quantity"`
}

func GetInventory(db *sql.DB, employeeId, merchId int64) (*Inventory, error) {
	env := "tables.inventory.GetInventory"

	stmt, err := db.Prepare("SELECT * FROM inventory WHERE employee_id = $1 AND merch_id = $2;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var invId int64
	var invEId int64
	var invMId int64
	var quantity int

	err = stmt.QueryRow(employeeId, merchId).Scan(&invId, &invEId, &invMId, &quantity)
	if err != nil {
		log.Printf("%s: failed to get the record, err: %v", env, err)
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Inventory = Inventory{Id: invId, EmployeeId: invEId, MerchId: invMId, Quantity: quantity}

	return &res, nil
}

func IncreaseQuantity(db *sql.DB, id int64, additionalQuantity int) error {
	env := "tables.inventory.IncreaseQuantity"

	getStmt, err := db.Prepare("SELECT quantity FROM inventory WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the select stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the select stmt, err: %w", env, err)
	}

	var quantity int
	err = getStmt.QueryRow(id).Scan(&quantity)
	if err != nil {
		log.Printf("%s: failed to get the record, err: %v", env, err)
		return fmt.Errorf("%s: %w", env, err)
	}

	quantity += additionalQuantity
	updStmt, err := db.Prepare("UPDATE inventory SET quantity = $2 WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the update stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the update stmt, err: %w", env, err)
	}

	_, err = updStmt.Exec(id, quantity)
	if err != nil {
		log.Printf("%s: failed to execute the update stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to execute the update stmt, err: %w", env, err)
	}

	return nil
}

func CreateInventoryRecord(db *sql.DB, inventory *Inventory) error {
	env := "tables.inventory.CreateInventoryRecord"

	stmt, err := db.Prepare("INSERT INTO inventory(employee_id, merch_id, quantity) VALUES($1, $2, $3);")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	_, err = stmt.Exec(inventory.EmployeeId, inventory.MerchId, inventory.Quantity)
	if err != nil {
		log.Printf("%s: unmatched arguments to insert, err: %v", env, err)
		return fmt.Errorf("%s: unmatched arguments to insert, err: %w", env, err)
	}

	return nil
}

func BuyInventory(db *sql.DB, employeeId, merchId int64, quantity int) error {
	env := "tables.inventory.GetInventory"

	record, err := GetInventory(db, employeeId, merchId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			inv := Inventory{EmployeeId: employeeId, MerchId: merchId, Quantity: quantity}
			createErr := CreateInventoryRecord(db, &inv)
			if createErr != nil {
				log.Printf("%s: failed to create the new inv record, err: %v", env, createErr)
				return fmt.Errorf("%s: failed to create the new inv record, err: %w", env, createErr)
			}
		} else {
			log.Printf("%s: failed to get the inv record, err: %v", env, err)
			return fmt.Errorf("%s: failed to get the inv record, err: %w", env, err)
		}
	}

	err = IncreaseQuantity(db, record.Id, quantity)
	if err != nil {
		log.Printf("%s: failed to exec the IncreaseQuantity func, err: %v", env, err)
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}
