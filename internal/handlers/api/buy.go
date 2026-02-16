package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/accounts"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/merch"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/api"
)

func BuyItemHandler(db *sql.DB, merchId int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := api.GetUser(r, db)
		if err != nil {
			log.Printf("cannot get the user by token, err: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
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

		if err := merch.BuyMerch(db, merchId, user.Id, price); err != nil {
			log.Printf("failed to buy the merch, err: %v", err)
			http.Error(w, "Failed to buy the merch", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
