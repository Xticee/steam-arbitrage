package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

type Item struct {
	Name         string    `json:"name"`
	MarketPrice  float64   `json:"market_price"`
	SteamPrice   float64   `json:"steam_price"`
	SteamBalance float64   `json:"steam_balance"`
	Profit       float64   `json:"profit"`
	IconId       string    `json:"icon_id"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewDb(connStr string) (*DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (d *DB) SaveMarketPrice(name string, marketPrice float64) error {
	query := `
	INSERT INTO items (name, market_price)
	VALUES ($1, $2)
	ON CONFLICT (name) DO UPDATE SET
		market_price = $2,
		profit = items.steam_price * 0.87 - $2 * 1.05821,
		updated_at = NOW()`
	_, err := d.conn.Exec(query, name, marketPrice)

	return err
}

func (d *DB) SaveSteamPrice(name string, steamPrice float64, iconId string) error {
	query := `
	INSERT INTO items (name, steam_price, icon_id)
	VALUES ($1, $2, $3)
	ON CONFLICT (name) DO UPDATE SET
		steam_price = $2,
		steam_balance = $2 * 0.87,
		profit = $2 * 0.87 - items.market_price * 1.05821,
		icon_id = $3,
		updated_at = NOW()`
	_, err := d.conn.Exec(query, name, steamPrice, iconId)

	return err
}

func (d *DB) GetTopItems(amount float64) ([]Item, error) {
	query := `
	SELECT * FROM items
	WHERE steam_balance < $1 AND profit > 0
	ORDER BY steam_balance DESC
	LIMIT 10`
	rows, err := d.conn.Query(query, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]Item, 0, 10)
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.Name, &item.MarketPrice, &item.SteamPrice, &item.SteamBalance, &item.Profit, &item.IconId, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}
