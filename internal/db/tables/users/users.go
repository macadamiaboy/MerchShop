package users

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func GetUserByLogin(db *sql.DB, login string) (*User, error) {
	env := "tables.employees.GetUserByLogin"

	stmt, err := db.Prepare("SELECT * FROM users WHERE login = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var userId int64
	var userLogin string
	var userPassword string

	err = stmt.QueryRow(login).Scan(&userId, &userLogin, &userPassword)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res User = User{Id: userId, Login: userLogin, Password: userPassword}

	return &res, nil
}

func CreateUser(db *sql.DB, user *User) error {
	env := "tables.employees.CreateUser"

	stmt, err := db.Prepare("INSERT INTO users(login, password) VALUES($1, $2);")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	_, err = stmt.Exec(user.Login, user.Password)
	if err != nil {
		log.Printf("%s: unmatched arguments to insert, err: %v", env, err)
		return fmt.Errorf("%s: unmatched arguments to insert, err: %w", env, err)
	}

	return nil
}
