package pokeapi

import "net/http"

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return NewClientWithHTTPClient(nil)
}

func NewClientWithHTTPClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}
