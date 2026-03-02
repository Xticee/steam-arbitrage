package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
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
	INSERT INTO items (name, rust_tm_price)
	VALUES ($1, $2)
	ON CONFLICT (name) DO UPDATE SET
		rust_tm_price = $2,
		profit = items.steam_price * 0.87 - $2 * 1.05821,
		updated_at = NOW()`
	_, err := d.conn.Exec(query, name, marketPrice)

	return err
}

func (d *DB) SaveSteamPrice(name string, steamPrice float64) error {
	query := `
	INSERT INTO items (name, steam_price)
	VALUES ($1, $2)
	ON CONFLICT (name) DO UPDATE SET
		steam_price = $2,
		steam_balance = EXCLUDED.steam_price * 0.87,
		profit = $2 * 0.87 - items.rust_tm_price * 1.05821,
		updated_at = NOW()`
	_, err := d.conn.Exec(query, name, steamPrice)

	return err
}
