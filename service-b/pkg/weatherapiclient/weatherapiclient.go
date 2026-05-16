package weatherapiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type CurrentResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func New(apiKey string, client *http.Client) *Client {
	if client == nil {
		client = &http.Client{}
	}
	return &Client{
		baseURL: "https://api.weatherapi.com/v1",
		apiKey:  apiKey,
		client:  client,
	}
}

func (c *Client) GetTemperature(location string) (float64, error) {
	endpoint := fmt.Sprintf("%s/current.json", c.baseURL)

	params := url.Values{}
	params.Add("key", c.apiKey)
	params.Add("q", location)
	params.Add("aqi", "no")

	resp, err := c.client.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data CurrentResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
