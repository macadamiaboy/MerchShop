package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/auth"
)

func getToken(r *http.Request) (string, error) {
	env := "helpers.api.auth.GetToken"

	token := r.Header.Get("Authorization")

	if token == "" {
		return "", fmt.Errorf("%s: the token is not provided", env)
	} else {
		return token, nil
	}
}

func GetUserId(r *http.Request, db *sql.DB) (int64, error) {
	env := "helpers.api.auth.GetUserId"

	token, err := getToken(r)
	if err != nil {
		log.Printf("no token provided in header, err: %v", err)
		return 0, fmt.Errorf("%s: no token provided in header, err: %w", env, err)
	}

	id, err := auth.GetIdFromToken(token)
	if err != nil {
		log.Printf("provided token is incorrect or has been expired, err: %v", err)
		return 0, fmt.Errorf("%s: provided token is incorrect or has been expired, err: %w", env, err)
	}

	return id, nil
}
