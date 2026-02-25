package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// LocationArea represents a named location area returned by the PokeAPI.
type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// LocationAreasResponse is the paginated response for the location-area
// list endpoint.
type LocationAreasResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

const locationAreasPageSize = 20

// GetLocationAreas fetches a page of location areas starting at the given
// offset.
func (c *Client) GetLocationAreas(offset int) (LocationAreasResponse, error) {
	if c == nil || c.httpClient == nil {
		return LocationAreasResponse{}, errors.New("pokeapi client is not initialized")
	}

	endpoint, err := url.Parse(c.baseURL + "/location-area")
	if err != nil {
		return LocationAreasResponse{}, err
	}

	query := endpoint.Query()
	query.Set("offset", strconv.Itoa(offset))
	query.Set("limit", strconv.Itoa(locationAreasPageSize))
	endpoint.RawQuery = query.Encode()

	data, found := c.cache.Get(endpoint.String())
	if !found {
		req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
		if err != nil {
			return LocationAreasResponse{}, err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return LocationAreasResponse{}, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationAreasResponse{}, err
		}

		if res.StatusCode > 299 {
			err = fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(data))
			return LocationAreasResponse{}, err
		}

		c.cache.Add(endpoint.String(), data)
	}

	locationAreasResponse := LocationAreasResponse{}
	err = json.Unmarshal(data, &locationAreasResponse)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return locationAreasResponse, nil
}
