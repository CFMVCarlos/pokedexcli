package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		userText := scanner.Text()
		cleanText := cleanInput(userText)
		if len(cleanText) == 0 {
			continue
		}

		commands := getCommands()
		cmd, ok := commands[cleanText[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := cmd.callback()
		if err != nil {
			fmt.Println(err)
		}
	}
}
