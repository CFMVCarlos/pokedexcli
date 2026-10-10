package cli

import (
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/repository"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(config *Config, args ...string) error
}

type Config struct {
	Commands map[string]CliCommand
	Next     *string
	Previous *string
	Client   *pokeapi.Client
	Pokedex  *repository.Pokedex
}
