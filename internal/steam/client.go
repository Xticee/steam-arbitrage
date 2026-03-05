package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const endpointPattern = "https://steamcommunity.com/market/search/render/?appid=252490&norender=1&currency=5&sort_column=name&sort_dir=asc&start="

type Client struct {
	http *http.Client
}

type steamItem struct {
	Name   string
	Price  float64
	IconID string
}

func NewClient(http *http.Client) *Client {
	return &Client{http: http}
}

func (c *Client) FetchSteamPrices(start int) ([]steamItem, error) {
	var raw struct {
		Result []steamItem `json:"results"`
	}

	startStr := strconv.Itoa(start)
	endpoint := fmt.Sprintf("%s%s", endpointPattern, startStr)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", os.Getenv("STEAM_COOKIES"))
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wrong responce: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return raw.Result, nil
}

func (i *steamItem) UnmarshalJSON(data []byte) error {
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
