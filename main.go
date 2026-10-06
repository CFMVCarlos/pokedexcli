package main

import (
	"fmt"
)

func main() {
	fmt.Println("Welcome to the Pokedex!")

	config := &config{
		commands: GetCommands(),
		next:     "https://pokeapi.co/api/v2/location-area",
		previous: "",
	}
	StartRepl(config)
}
