package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartRepl(config *config) {
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

		commands := config.commands
		cmd, ok := commands[cleanText[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := cmd.callback(config)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	cleanText := strings.TrimSpace(text)
	words := strings.Fields(strings.ToLower(cleanText))
	return words
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range config.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func GetCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	return commands
}
