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
	name  string
	price float64
}

func NewClient(http *http.Client) *Client {
	return &Client{http: http}
}

func (c *Client) FetchSteamPrices(start int) (map[string]float64, error) {
	var itemsMap = make(map[string]float64)
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

	for _, item := range raw.Result {
		itemsMap[item.name] = item.price
	}

	return itemsMap, nil
}

func (i *steamItem) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name  string `json:"name"`
		Price int    `json:"sell_price"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	floatPrice := float64(raw.Price) / 100

	i.name = raw.Name
	i.price = floatPrice

	return nil
}
