package main

import (
	"log"
	"net/http"
	"os"
	"parser/internal/api"
	"parser/internal/db"
	"parser/internal/market"
	"parser/internal/steam"
	"parser/internal/updater"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	marketClient := market.NewClient(client)
	steamClient := steam.NewClient(client, os.Getenv("STEAM_COOKIES"))

	itemsBase, err := db.NewDb(os.Getenv("DB_CONN"))
	if err != nil {
		log.Fatal(err)
	}

	updater.Run(itemsBase, marketClient, steamClient)

	router := chi.NewRouter()
	handler := api.NewHandler(itemsBase)
	router.Get("/", handler.GeneralPage)
	router.Get("/search", handler.GetItems)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
