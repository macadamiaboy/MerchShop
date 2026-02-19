package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/accounts"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/inventory"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db/tables/transfers"
	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/api"
)

type infoResponse struct {
	Coins       int                    `json:"coins"`
	Inventory   *[]inventory.Inv       `json:"inventory"`
	CoinHistory *transfers.CoinHistory `json:"coinHistory"`
}

func InfoHandler(db *sql.DB, merchId int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := api.GetUser(r, db)
		if err != nil {
			log.Printf("cannot get the user by token, err: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		balance, err := accounts.GetBalanceById(db, user.Id)
		if err != nil {
			log.Printf("failed to get the balance, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		inventory, err := inventory.GetAllUsersInventory(db, user.Id)
		if err != nil {
			log.Printf("failed to get the inventory, err: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		coinHistory, err := transfers.GetCoinHistory(db, user.Id)
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
