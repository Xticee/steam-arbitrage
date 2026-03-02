package rusttm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const apiURL = "https://rust.tm/api/v2/prices/RUB.json"

type Client struct {
	http *http.Client
}

type apiResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name  string
	Price float64
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

	var jsonResp apiResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return nil, err
	}

	return jsonResp.Items, nil
}

func (i *Item) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name  string `json:"market_hash_name"`
		Price string `json:"price"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	i.Name = raw.Name
	price, err := strconv.ParseFloat(raw.Price, 32)
	if err != nil {
		return err
	}
	i.Price = float64(price)

	return nil
}
