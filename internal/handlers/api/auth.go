package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/accounts"
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

		//log.Println("basic job is done")

		var tokenId int64

		curUser, getErr := users.GetUserByLogin(db, requestBody.Username)
		if getErr != nil {
			if errors.Is(getErr, sql.ErrNoRows) {
				log.Println("creating user")

				user := users.User{
					Login:    requestBody.Username,
					Password: password,
				}

				userId, createErr := users.CreateUser(db, &user)
				if createErr != nil {
					log.Printf("failed to create the new user, err: %v", createErr)
					http.Error(w, createErr.Error(), http.StatusInternalServerError)
					return
				}

				// save the created user's id
				tokenId = userId

				accErr := accounts.CreateAccount(db, userId)
				if accErr != nil {
					log.Printf("failed to create the new account, err: %v", createErr)
					http.Error(w, accErr.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				log.Printf("failed to get the user, err: %v", getErr)
				http.Error(w, getErr.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			//log.Println("checking the password")
			if correctPassword := hash.CheckPasswordHash(requestBody.Password, curUser.Password); !correctPassword {
				log.Printf("incorrect password")
				http.Error(w, "Incorrect login or passsword", http.StatusUnauthorized)
				return
			}

			// save the found user's id
			tokenId = curUser.Id
		}

		//log.Println("generating the token")

		// token id is either the id of the user who was created right now or was found in the db by login
		token, err := auth.GenToken(tokenId)
		if err != nil {
			log.Printf("failed to create the token, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//log.Println("token was generated")

		response := AuthResponse{
			Token: token,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
