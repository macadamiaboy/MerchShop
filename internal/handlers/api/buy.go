package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/accounts"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/merch"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/users"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/api"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/auth"
)

func BuyItemHandler(db *sql.DB, merchId int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := api.GetToken(r)
		if err != nil {
			log.Printf("no token provided in header, err: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		login, err := auth.GetLoginFromToken(token)
		if err != nil {
			log.Printf("provided token is incorrect or has been expired, err: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := users.GetUserByLogin(db, login)
		if err != nil {
			log.Printf("cannot get the user by login, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		price, err := merch.GetMerchPrice(db, merchId)
		if err != nil {
			log.Printf("failed to get the merch price, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//is it needed to multiply by the quantity

		balance, err := accounts.GetBalanceById(db, user.Id)
		if err != nil {
			log.Printf("failed to get the balance, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if balance < price {
			log.Printf("not enough funds, err: %v", err)
			http.Error(w, "There's no enough funds on your account", http.StatusBadRequest)
			return
		}

		//transaction

		//
		//
		//

		response := AuthResponse{
			Token: token,
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
