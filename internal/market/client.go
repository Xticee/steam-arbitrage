package market

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiURL = "https://lis-skins.com/market_export_json/rust.json"
const exchangeRateUrl = "https://api.exchangerate-api.com/v4/latest/USD"

type Client struct {
	http *http.Client
}

type Item struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func NewClient(http *http.Client) *Client {
	return &Client{http: http}
}

func (c *Client) FetchItems() ([]Item, error) {
	resp, err := c.http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result := make([]Item, 0, 5000)
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	rate, err := c.fetchUSDRate()
	if err != nil {
		return nil, err
	}
	for i := range result {
		result[i].Price *= rate
	}

	return result, nil
}

func (c *Client) fetchUSDRate() (float64, error) {
	var rub struct {
		Rates struct {
			Rub float64 `json:"RUB"`
		} `json:"rates"`
	}
	resp, err := c.http.Get(exchangeRateUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&rub)

	return rub.Rates.Rub, err
}
