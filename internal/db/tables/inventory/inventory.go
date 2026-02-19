package inventory

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/merch"
)

type Inventory struct {
	Id         int64 `json:"id"`
	EmployeeId int64 `json:"employee_id"`
	MerchId    int64 `json:"merch_id"`
	Quantity   int   `json:"quantity"`
}

type Inv struct {
	InvType  string `json:"type"`
	Quantity int    `json:"quantity"`
}

func GetInventory(tx *sql.Tx, employeeId, merchId int64) (*Inventory, error) {
	env := "tables.inventory.GetInventory"

	stmt, err := tx.Prepare("SELECT * FROM inventory WHERE employee_id = $1 AND merch_id = $2;")
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

func IncreaseQuantity(tx *sql.Tx, id int64, additionalQuantity int) error {
	env := "tables.inventory.IncreaseQuantity"

	getStmt, err := tx.Prepare("SELECT quantity FROM inventory WHERE id = $1;")
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
	updStmt, err := tx.Prepare("UPDATE inventory SET quantity = $2 WHERE id = $1;")
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

func CreateInventoryRecord(tx *sql.Tx, inventory *Inventory) error {
	env := "tables.inventory.CreateInventoryRecord"

	stmt, err := tx.Prepare("INSERT INTO inventory(employee_id, merch_id, quantity) VALUES($1, $2, $3);")
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

func (inv *Inventory) BuyInventory(tx *sql.Tx, employeeId, merchId int64 /*, quantity int*/) error {
	env := "tables.inventory.BuyInventory"

	//is prepared make purchases in multiple copies
	quantity := 1

	record, err := GetInventory(tx, employeeId, merchId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			//if there's no record with such user and merch, create:
			inv := Inventory{EmployeeId: employeeId, MerchId: merchId, Quantity: quantity}
			createErr := CreateInventoryRecord(tx, &inv)
			if createErr != nil {
				log.Printf("%s: failed to create the new inv record, err: %v", env, createErr)
				return fmt.Errorf("%s: failed to create the new inv record, err: %w", env, createErr)
			}
			return nil
		} else {
			log.Printf("%s: failed to get the inv record, err: %v", env, err)
			return fmt.Errorf("%s: failed to get the inv record, err: %w", env, err)
		}
	}

	//if there is such record, increase the amount of user's merch of this type
	err = IncreaseQuantity(tx, record.Id, quantity)
	if err != nil {
		log.Printf("%s: failed to exec the IncreaseQuantity func, err: %v", env, err)
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func GetAllUsersInventory(db *sql.DB, employeeId int64) (*[]Inv, error) {
	env := "tables.inventory.GetAllUsersInventory"

	rows, err := db.Query("SELECT id, quantity FROM inventory WHERE employee_id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}
	defer rows.Close()

	var collection []Inv
	for rows.Next() {
		var inventory Inventory
		if err := rows.Scan(&inventory.Id, &inventory.Quantity); err != nil {
			log.Printf("%s: failed to get the inventory record, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the inventory record, err: %w", env, err)
		}

		name, err := merch.GetMerchName(db, inventory.Id)
		if err != nil {
			log.Printf("%s: failed to get the merch type, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the merch type, err: %w", env, err)
		}

		collection = append(collection, Inv{InvType: name, Quantity: inventory.Quantity})
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}
