package main

import (
	"fmt"
	"time"

	"pokedexcli/internal/cli"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/repository"
)

func main() {
	fmt.Println("Welcome to the Pokedex!")
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	pokedexRepo := repository.NewPokedex("pokedex.json")

	config := &cli.Config{
		Commands: cli.GetCommands(),
		Client:   &pokeClient,
		Pokedex:  pokedexRepo,
	}

	cli.StartRepl(config)
}
