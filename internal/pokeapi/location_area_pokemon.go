package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// LocationAreaResponse is the response for a single location area,
// containing the Pokemon that can be encountered there.
type LocationAreaResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// Pokemon represents a Pokemon species identified by name.
type Pokemon struct {
	Name string `json:"name"`
}

// PokemonEncounter pairs a Pokemon with the location area it appears in.
type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

// GetLocationAreaPokemon fetches the Pokemon encounters for the named
// location area.
func (c *Client) GetLocationAreaPokemon(locationAreaName string) (LocationAreaResponse, error) {
	if c == nil || c.httpClient == nil {
		return LocationAreaResponse{}, errors.New("pokeapi client is not initialized")
	}

	endpoint := c.baseURL + "/location-area/" + locationAreaName

	data, found := c.cache.Get(endpoint)
	if !found {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return LocationAreaResponse{}, err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return LocationAreaResponse{}, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationAreaResponse{}, err
		}

		if res.StatusCode > 299 {
			err = fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(data))
			return LocationAreaResponse{}, err
		}

		c.cache.Add(endpoint, data)
	}

	var response LocationAreaResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	return response, nil
}
