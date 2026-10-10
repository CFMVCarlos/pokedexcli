package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pokedexcli/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache      pokecache.Cache
}

func NewClient(timeout time.Duration, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
}

func (c *Client) GetLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area"
	if pageURL != nil && *pageURL != "" {
		url = *pageURL
	}
	var data LocationAreaResponse
	err := c.fetch(url, &data)
	return data, err
}

func (c *Client) GetLocationArea(locationName string) (LocationArea, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationName)
	var data LocationArea
	err := c.fetch(url, &data)
	return data, err
}

func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	var data Pokemon
	err := c.fetch(url, &data)
	return data, err
}

func (c *Client) fetch(url string, target interface{}) error {
	if val, ok := c.cache.Get(url); ok {
		return json.Unmarshal(val, target)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.cache.Add(url, body)
	return json.Unmarshal(body, target)
}
