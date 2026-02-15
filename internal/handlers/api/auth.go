package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/users"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/auth"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/hash"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	Errors string `json:"errors"`
}

func AuthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody *AuthRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			log.Printf("failed to get the request body, err: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		password, err := hash.HashPassword(requestBody.Password)
		if err != nil {
			log.Printf("failed to generate the hash, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := users.User{
			Login:    requestBody.Username,
			Password: password,
		}

		curUser, getErr := users.GetUserByLogin(db, user.Login)
		if getErr != nil {
			if errors.Is(getErr, sql.ErrNoRows) {
				createErr := users.CreateUser(db, &user)
				if createErr != nil {
					log.Printf("failed to create the new user, err: %v", createErr)
					http.Error(w, createErr.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				log.Printf("failed to get the user, err: %v", getErr)
				http.Error(w, getErr.Error(), http.StatusInternalServerError)
				return
			}
		}

		if correctPassword := hash.CheckPasswordHash(requestBody.Password, curUser.Password); !correctPassword {
			log.Printf("incorrect password")
			http.Error(w, "Incorrect login or passsword", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenToken(user.Login)
		if err != nil {
			log.Printf("failed to create the token, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := AuthResponse{
			Token: token,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
