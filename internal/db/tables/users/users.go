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
	env := "tables.users.GetUserByLogin"

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

func CreateUser(db *sql.DB, user *User) (int64, error) {
	env := "tables.users.CreateUser"

	stmt, err := db.Prepare("INSERT INTO users(login, password) VALUES($1, $2) RETURNING id;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return 0, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var userId int64
	err = stmt.QueryRow(user.Login, user.Password).Scan(&userId)
	if err != nil {
		log.Printf("%s: unmatched arguments to insert, err: %v", env, err)
		return 0, fmt.Errorf("%s: unmatched arguments to insert, err: %w", env, err)
	}

	return userId, nil
}
