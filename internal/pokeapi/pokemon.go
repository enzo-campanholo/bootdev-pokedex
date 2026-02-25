package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Type struct {
	Name string `json:"name"`
}

type Stat struct {
	Name string `json:"name"`
}

type PokemonType struct {
	Type Type `json:"type"`
}

type PokemonStat struct {
	Stat     Stat `json:"stat"`
	BaseStat int  `json:"base_stat"`
}

// Pokemon represents a Pokemon species identified by name.
type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	if c == nil || c.httpClient == nil {
		return Pokemon{}, errors.New("pokeapi client is not initialized")
	}

	endpoint := c.baseURL + "/pokemon/" + pokemonName

	data, found := c.cache.Get(endpoint)
	if !found {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return Pokemon{}, err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return Pokemon{}, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, err
		}

		if res.StatusCode > 299 {
			err = fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(data))
			return Pokemon{}, err
		}

		c.cache.Add(endpoint, data)
	}

	var pokemon Pokemon
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
