package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/users"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/auth"
)

func GetToken(r *http.Request) (string, error) {
	env := "helpers.api.auth.GetToken"

	token := r.Header.Get("Authorization")

	if token == "" {
		return "", fmt.Errorf("%s: the token is not provided", env)
	} else {
		return token, nil
	}
}

func GetUser(r *http.Request, db *sql.DB) (*users.User, error) {
	env := "helpers.api.auth.GetUser"

	token, err := GetToken(r)
	if err != nil {
		log.Printf("no token provided in header, err: %v", err)
		return nil, fmt.Errorf("%s: no token provided in header, err: %w", env, err)
	}

	login, err := auth.GetLoginFromToken(token)
	if err != nil {
		log.Printf("provided token is incorrect or has been expired, err: %v", err)
		return nil, fmt.Errorf("%s: provided token is incorrect or has been expired, err: %w", env, err)
	}

	user, err := users.GetUserByLogin(db, login)
	if err != nil {
		log.Printf("cannot get the user by login, err: %v", err)
		return nil, fmt.Errorf("%s: cannot get the user by login, err: %w", env, err)
	}

	return user, nil
}
