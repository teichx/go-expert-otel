package viacepclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type Response struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

func New(client *http.Client) *Client {
	if client == nil {
		client = &http.Client{}
	}
	return &Client{
		baseURL: "https://viacep.com.br/ws",
		client:  client,
	}
}

func (c *Client) GetLocation(zipcode string) (string, error) {
	url := fmt.Sprintf("%s/%s/json", c.baseURL, zipcode)

	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	if data.Erro {
		return "", fmt.Errorf("zipcode not found")
	}

	return data.Localidade, nil
}
