package updater

import (
	"log"
	"parser/internal/db"
	"parser/internal/market"
	"parser/internal/steam"
	"time"
)

const totalSteamItems = 5120

func Run(db *db.DB, marketClient *market.Client, steamClient *steam.Client) {
	go updateMarketData(marketClient, db)
	go updateSteamData(steamClient, db)
}

func updateMarketData(marketClient *market.Client, db *db.DB) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		items, err := marketClient.FetchItems()
		if err != nil {
			log.Printf("Ошибка при попытке парсинга цен с маркета: %v", err)
			continue
		}
		for _, item := range items {
			if err = db.SaveMarketPrice(item.Name, item.Price); err != nil {
				log.Printf("Ошибка при попытке сохранения цены маркета в базу: %v", err)
			}
		}

		<-ticker.C
	}
}

func updateSteamData(steamClient *steam.Client, db *db.DB) {
	ticker := time.NewTicker(5*time.Minute + 2*time.Second)
	defer ticker.Stop()

	var start int
	for {
		for i := 0; i < 25; i++ {
			itemSlice, err := steamClient.FetchSteamPrices(start)
			if err != nil {
				log.Printf("Ошибка при попытке парсинга steam цен: %v", err)
				continue
			}
			for _, item := range itemSlice {
				if err = db.SaveSteamPrice(item.Name, item.Price, item.IconID); err != nil {
					log.Printf("Ошибка при попытке записи в базу steam данных: %v", err)
				}
			}
			start += 10
			if start == totalSteamItems {
				start = 0
			}
			time.Sleep(1 * time.Second)
		}
		<-ticker.C
	}
}
