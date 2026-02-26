package pokeapi

// LocationArea represents a named location area returned by the PokeAPI.
type LocationArea struct {
	Name string `json:"name"`
}

// GetLocationAreas fetches a page of location areas. Pass an empty pageURL
// to fetch the first page. Returns the areas and the next/previous page URLs
// (empty string means no more pages in that direction).
func (c *Client) GetLocationAreas(pageURL string) ([]LocationArea, string, string, error) {
	if pageURL == "" {
		pageURL = c.baseURL + "/location-area"
	}

	var resp struct {
		Next     *string        `json:"next"`
		Previous *string        `json:"previous"`
		Results  []LocationArea `json:"results"`
	}
	if err := c.get(pageURL, &resp); err != nil {
		return nil, "", "", err
	}

	next := ""
	if resp.Next != nil {
		next = *resp.Next
	}
	prev := ""
	if resp.Previous != nil {
		prev = *resp.Previous
	}

	return resp.Results, next, prev, nil
}
