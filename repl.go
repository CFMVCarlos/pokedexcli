package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

// StartRepl initiates the interactive read-eval-print loop (REPL), scanning standard input
// for user commands, parsing arguments, and executing registered callbacks.
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
		args := []string{}
		if len(cleanText) > 1 {
			args = cleanText[1:]
		}
		err := cmd.callback(config, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// cleanInput trims outer whitespace and splits the input string into lowercase word tokens.
func cleanInput(text string) []string {
	cleanText := strings.TrimSpace(text)
	words := strings.Fields(strings.ToLower(cleanText))
	return words
}

// commandExit terminates the Pokedex CLI application cleanly with exit code 0.
func commandExit(config *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

// commandHelp displays usage instructions and lists all registered commands with their descriptions.
func commandHelp(config *config, args ...string) error {
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

// commandMap fetches and displays the next 20 location areas from the PokeAPI.
func commandMap(config *config, args ...string) error {
	if config.next == "" {
		return fmt.Errorf("you're on the last page")
	}

	return printLocationAreas(config, config.next)
}

// commandMapb fetches and displays the previous 20 location areas from the PokeAPI.
func commandMapb(config *config, args ...string) error {
	if config.previous == "" {
		return fmt.Errorf("you're on the first page")
	}

	return printLocationAreas(config, config.previous)
}

// printLocationAreas retrieves location areas from the given URL (via cache or HTTP),
// prints the names of the areas, and updates next/previous pagination state.
func printLocationAreas(config *config, url string) error {
	var body []byte

	if val, ok := config.cache.Get(url); ok {
		body = val
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		body = data
		config.cache.Add(url, body)
	}

	var data locationAreaResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	if data.Next != nil {
		config.next = *data.Next
	} else {
		config.next = ""
	}

	if data.Previous != nil {
		config.previous = *data.Previous
	} else {
		config.previous = ""
	}

	return nil
}

// commandExplore lists all Pokémon that can be encountered within a specified location area.
func commandExplore(config *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location name")
	}
	location := args[0]
	fmt.Printf("Exploring %s...\n", location)

	fullURL := "https://pokeapi.co/api/v2/location-area/" + location
	var body []byte

	if val, ok := config.cache.Get(fullURL); ok {
		body = val
	} else {
		req, err := http.Get(fullURL)
		if err != nil {
			return err
		}
		defer req.Body.Close()

		data, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		body = data
		config.cache.Add(fullURL, data)
	}

	var data LocationArea
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemonEncounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

// commandCatch attempts to catch a specified Pokémon, calculating catch probability based
// on base experience and adding caught Pokémon to the user's Pokédex.
func commandCatch(config *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	pokemonName := strings.ToLower(args[0])
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	fullURL := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	var body []byte

	if val, ok := config.cache.Get(fullURL); ok {
		body = val
	} else {
		req, err := http.Get(fullURL)
		if err != nil {
			return err
		}
		defer req.Body.Close()

		if req.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to find pokemon: %s", req.Status)
		}

		data, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		body = data
		config.cache.Add(fullURL, data)
	}

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return err
	}

	baseExp := pokemon.BaseExperience
	if baseExp <= 0 {
		baseExp = 1
	}

	res := rand.Intn(baseExp)
	if res < 40 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		if config.pokedex == nil {
			config.pokedex = make(map[string]Pokemon)
		}
		config.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

// commandInspect prints comprehensive stats, height, weight, and elemental types
// for a previously captured Pokémon.
func commandInspect(config *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	pokemonName := strings.ToLower(args[0])

	pokemon, exists := config.pokedex[pokemonName]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

// commandPokedex lists all Pokémon names that have been captured and saved in the user's Pokédex.
func commandPokedex(config *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

// GetCommands creates and returns the registry mapping command names to their cliCommand definitions.
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

	commands["map"] = cliCommand{
		name:        "map",
		description: "Displays the next 20 locations",
		callback:    commandMap,
	}

	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays the previous 20 locations",
		callback:    commandMapb,
	}

	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Shows all pokemon in a given area",
		callback:    commandExplore,
	}

	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Attempt to catch a pokemon",
		callback:    commandCatch,
	}

	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "View details about a caught pokemon",
		callback:    commandInspect,
	}

	commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "List all caught pokemon",
		callback:    commandPokedex,
	}

	return commands
}
