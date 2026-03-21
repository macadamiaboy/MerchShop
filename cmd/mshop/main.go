package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/macadamiaboy/AvitoMerchShop/internal/config"
	"github.com/macadamiaboy/AvitoMerchShop/internal/db"
	handlers "github.com/macadamiaboy/AvitoMerchShop/internal/handlers/api"
	localMW "github.com/macadamiaboy/AvitoMerchShop/internal/middleware"
)

func main() {

	//config
	cfg := config.LoadConfigData()
	addr := fmt.Sprintf("%s:%v", cfg.Server.Host, cfg.Server.Port)

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	//dbinit
	if err := db.Init(); err != nil {
		fmt.Println(err)
		log.Fatal("failed to create the db")
	}

	// opening db
	db, err := db.PrepareDB()
	if err != nil {
		log.Fatalf("failed to prepare the db: %v", err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
			return
		}
	}()

	router.Route("/api", func(r chi.Router) {
		r.Post("/auth", handlers.AuthHandler(db.Connection))

		r.Group(func(r chi.Router) {
			r.Use(localMW.AuthMiddleware)

			r.Get("/info", handlers.InfoHandler(db.Connection))
			r.Post("/sendCoin", handlers.SendCoinHandler(db.Connection))
			r.Get("/but/{item}", handlers.BuyItemHandler(db.Connection))
		})
	})

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	log.Printf("starting server. address: %s", addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("failed to start the server")
	}

	log.Fatal("server stoppped")

}
