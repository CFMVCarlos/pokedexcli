package main

import "pokedexcli/internal/pokecache"

type cliCommand struct {
	name        string
	description string
	callback    func(config *config) error
}
type config struct {
	commands map[string]cliCommand
	next     string
	previous string
	cache    pokecache.Cache
}

type locationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
