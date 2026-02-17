package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/accounts"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/users"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/api"
)

type RequestBody struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func SendCoinHandler(db *sql.DB, merchId int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := api.GetUser(r, db)
		if err != nil {
			log.Printf("cannot get the user by token, err: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var requestBody *RequestBody

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			log.Printf("failed to get the request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		receiver, err := users.GetUserByLogin(db, requestBody.ToUser)
		if err != nil {
			log.Printf("there's no such user with the provided login, err: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		balance, err := accounts.GetBalanceById(db, user.Id)
		if err != nil {
			log.Printf("failed to get the balance, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if balance < requestBody.Amount {
			log.Printf("not enough funds, err: %v", err)
			http.Error(w, "There's no enough funds on your account", http.StatusBadRequest)
			return
		}

		if err := accounts.Transfer(db, user.Id, receiver.Id, requestBody.Amount); err != nil {
			log.Printf("failed to buy the merch, err: %v", err)
			http.Error(w, "Failed to buy the merch", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
