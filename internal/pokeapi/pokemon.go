package pokeapi

import "encoding/json"

// Stat holds a single base stat for a Pokemon.
type Stat struct {
	Name     string
	BaseStat int
}

// Pokemon represents a Pokemon species identified by name.
type Pokemon struct {
	Name           string
	BaseExperience int
	Height         int
	Weight         int
	Stats          []Stat
	Types          []string
}

func (p *Pokemon) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name           string `json:"name"`
		BaseExperience int    `json:"base_experience"`
		Height         int    `json:"height"`
		Weight         int    `json:"weight"`
		Stats          []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.Name = raw.Name
	p.BaseExperience = raw.BaseExperience
	p.Height = raw.Height
	p.Weight = raw.Weight

	p.Stats = make([]Stat, len(raw.Stats))
	for i, s := range raw.Stats {
		p.Stats[i] = Stat{Name: s.Stat.Name, BaseStat: s.BaseStat}
	}

	p.Types = make([]string, len(raw.Types))
	for i, t := range raw.Types {
		p.Types[i] = t.Type.Name
	}

	return nil
}

// GetPokemon fetches a Pokemon by name, using the cache when available.
func (c *Client) GetPokemon(name string) (Pokemon, error) {
	var pokemon Pokemon
	err := c.get(c.baseURL+"/pokemon/"+name, &pokemon)
	return pokemon, err
}
