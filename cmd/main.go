package main

import (
	"log"
	"net/http"
	"os"
	"parser/internal/db"
	"parser/internal/rusttm"
	"parser/internal/steam"
	"parser/internal/updater"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	marketClient := rusttm.NewClient(client)
	steamClient := steam.NewClient(client)

	connStr := os.Getenv("DB_CONN")
	itemsBase, err := db.NewDb(connStr)
	if err != nil {
		log.Fatal(err)
	}

	updater.Run(itemsBase, marketClient, steamClient)
	if err != nil {
		log.Fatal(err)
	}
}
