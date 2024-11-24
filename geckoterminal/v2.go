package geckoterminal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var baseURL = "https://api.geckoterminal.com/api/v2"

// Client struct
type Client struct {
	httpClient *http.Client
}

type SimpleSinglePrice struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			TokenPrices map[string]string `json:"token_prices"`
		} `json:"attributes"`
	} `json:"data"`
}

// NewClient create new client object
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{httpClient: httpClient}
}

// helper
// doReq HTTP client
func doReq(req *http.Request, client *http.Client) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

// MakeReq HTTP request helper
func (c *Client) MakeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	resp, err := doReq(req, c.httpClient)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// APIS

// SimplePrice /simple/price Multiple ID and Currency (ids, vs_currencies)
/*

https://api.geckoterminal.com/api/v2/simple/networks/arbitrum_nova/token_price/0x8afe4055ebc86bd2afb3940c0095c9aca511d852%2C0x722e8bdd2ce80a4422e880164f2079488e115365

{
  "data": {
    "id": "83d901ba-b3fc-4200-91b4-c9a672dd1668",
    "type": "simple_token_price",
    "attributes": {
      "token_prices": {
        "0x722e8bdd2ce80a4422e880164f2079488e115365": "3252.69554120332",
        "0x8afe4055ebc86bd2afb3940c0095c9aca511d852": "36.8621806437013"
      }
    }
  }
}
*/
func (c *Client) SimpleTokenPrice(network string, tokenAddresses []string) (map[string]float64, error) {
	addresssesParam := url.QueryEscape(strings.Join(tokenAddresses, ","))

	url := fmt.Sprintf("%s/simple/networks/%s/token_price/%s", baseURL, network, addresssesParam)
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	var prices SimpleSinglePrice
	err = json.Unmarshal(resp, &prices)
	if err != nil {
		return nil, err
	}

	tokenPrices := make(map[string]float64)
	for k, v := range prices.Data.Attributes.TokenPrices {
		if price, err := strconv.ParseFloat(v, 64); err == nil {
			tokenPrices[k] = price
		}
	}

	return tokenPrices, nil
}
