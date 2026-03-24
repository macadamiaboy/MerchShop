package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/accounts"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/inventory"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/transfers"
	localMW "github.com/macadamiaboy/AvitoMerchShop/internal/middleware"
)

type infoResponse struct {
	Coins       int                    `json:"coins"`
	Inventory   *[]inventory.Inv       `json:"inventory"`
	CoinHistory *transfers.CoinHistory `json:"coinHistory"`
}

func InfoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId, ok := r.Context().Value(localMW.UserIdKey).(int64)
		if !ok {
			log.Printf("Cannot assign the user id from the context to int64")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		balance, err := accounts.GetBalanceById(db, userId)
		if err != nil {
			log.Printf("failed to get the balance, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		inventory, err := inventory.GetAllUsersInventory(db, userId)
		if err != nil {
			log.Printf("failed to get the inventory, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		coinHistory, err := transfers.GetCoinHistory(db, userId)
		if err != nil {
			log.Printf("failed to get the coin history, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := infoResponse{
			Coins:       balance,
			Inventory:   inventory,
			CoinHistory: coinHistory,
		}

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
