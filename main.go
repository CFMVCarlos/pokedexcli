// Package main serves as the entrypoint for the Pokedex CLI application.
package main

import (
	"fmt"
	"time"

	"pokedexcli/internal/pokecache"
)

// main initializes application state, the TTL cache, and starts the REPL.
func main() {
	fmt.Println("Welcome to the Pokedex!")

	config := &config{
		commands: GetCommands(),
		next:     "https://pokeapi.co/api/v2/location-area",
		previous: "",
		cache:    pokecache.NewCache(5 * time.Second),
		pokedex:  make(map[string]Pokemon),
		saveFile: "pokedex.json",
	}
	StartRepl(config)
}
