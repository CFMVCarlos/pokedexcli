package main

import (
	"fmt"
	"time"

	"pokedexcli/internal/pokecache"
)

func main() {
	fmt.Println("Welcome to the Pokedex!")

	config := &config{
		commands: GetCommands(),
		next:     "https://pokeapi.co/api/v2/location-area",
		previous: "",
		cache:    pokecache.NewCache(5 * time.Second),
		pokedex:  make(map[string]Pokemon),
	}
	StartRepl(config)
}
