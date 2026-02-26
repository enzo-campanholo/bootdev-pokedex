// Package pokeapi provides a client for the PokeAPI REST service.
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

// Client is an HTTP client for the PokeAPI with built-in caching.
type Client struct {
	httpClient *http.Client
	baseURL    string
	cache      *pokecache.Cache
}

// NewClient creates a Client that uses a default http.Client and the
// provided cache.
func NewClient(cache *pokecache.Cache) *Client {
	return NewClientWithHTTPClient(nil, cache)
}

// NewClientWithHTTPClient creates a Client with the given http.Client. If
// httpClient is nil, a default http.Client is used.
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

func (c *Client) get(endpoint string, result any) error {
	data, found := c.cache.Get(endpoint)
	if !found {
		res, err := c.httpClient.Get(endpoint)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode > 299 {
			return fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(data))
		}

		c.cache.Add(endpoint, data)
	}
	return json.Unmarshal(data, result)
}
