package main

import (
	"fmt"
)

func main() {
	fmt.Println("Welcome to the Pokedex!")

	config := &config{
		commands: GetCommands(),
	}
	StartRepl(config)
}
