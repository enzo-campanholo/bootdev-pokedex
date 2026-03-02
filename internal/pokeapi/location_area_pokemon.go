package pokeapi

// GetLocationAreaPokemon returns the names of Pokemon found in the named
// location area.
func (c *Client) GetLocationAreaPokemon(name string) ([]string, error) {
	var resp struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}
	if err := c.get(c.baseURL+"/location-area/"+name, &resp); err != nil {
		return nil, err
	}

	names := make([]string, len(resp.PokemonEncounters))
	for i, enc := range resp.PokemonEncounters {
		names[i] = enc.Pokemon.Name
	}
	return names, nil
}
