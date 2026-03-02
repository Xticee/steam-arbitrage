package updater

import (
	"log"
	"parser/internal/db"
	"parser/internal/rusttm"
	"parser/internal/steam"
	"time"
)

const totalSteamItems = 5120

func Run(db *db.DB, marketClient *rusttm.Client, steamClient *steam.Client) {
	go updateMarketData(marketClient, db)
	updateSteamData(steamClient, db)
}

func updateMarketData(marketClient *rusttm.Client, db *db.DB) {
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
			itemMap, err := steamClient.FetchSteamPrices(start)
			if err != nil {
				log.Printf("Ошибка при попытке парсинга steam цен: %v", err)
				continue
			}
			for name, price := range itemMap {
				if err = db.SaveSteamPrice(name, price); err != nil {
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
