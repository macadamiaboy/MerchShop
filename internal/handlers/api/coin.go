package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/accounts"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/users"
	localMW "github.com/macadamiaboy/AvitoMerchShop/internal/middleware"
)

type sCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func SendCoinHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId, ok := r.Context().Value(localMW.UserIdKey).(int64)
		if !ok {
			log.Printf("Cannot assign the user id from the context to int64")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		var requestBody *sCoinRequest

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

		balance, err := accounts.GetBalanceById(db, userId)
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

		if err := accounts.Transfer(db, userId, receiver.Id, requestBody.Amount); err != nil {
			log.Printf("failed to send coins, err: %v", err)
			http.Error(w, "Failed to send coins", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
