package pokeapi

import (
	"net/http"

	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokecache"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	cache      *pokecache.Cache
}

func NewClient(cache *pokecache.Cache) *Client {
	return NewClientWithHTTPClient(nil, cache)
}

func NewClientWithHTTPClient(httpClient *http.Client, cache *pokecache.Cache) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		cache:      cache,
	}
}
