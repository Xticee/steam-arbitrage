package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseUrl = "https://steamcommunity.com/market/search/render/?appid=252490&norender=1&currency=5&sort_column=name&sort_dir=asc&start="

type Client struct {
	http    *http.Client
	cookies string
}

type SteamItem struct {
	Name   string
	Price  float64
	IconID string
}

func NewClient(http *http.Client, cookies string) *Client {
	return &Client{http: http, cookies: cookies}
}

func (c *Client) FetchSteamPrices(start int) ([]SteamItem, error) {
	var raw struct {
		Result []SteamItem `json:"results"`
	}

	endpoint := fmt.Sprintf("%s%d", baseUrl, start)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", c.cookies)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wrong response: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return raw.Result, nil
}

func (i *SteamItem) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name  string `json:"name"`
		Price int    `json:"sell_price"`
		Asset struct {
			IconId string `json:"icon_url"`
		} `json:"asset_description"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	floatPrice := float64(raw.Price) / 100

	i.Name = raw.Name
	i.Price = floatPrice
	i.IconID = raw.Asset.IconId

	return nil
}
