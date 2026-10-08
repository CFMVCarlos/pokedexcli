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

func cleanInput(text string) []string {
	cleanText := strings.TrimSpace(text)
	words := strings.Fields(strings.ToLower(cleanText))
	return words
}

func commandExit(config *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

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

func commandMap(config *config, args ...string) error {
	if config.next == "" {
		return fmt.Errorf("you're on the last page")
	}

	return printLocationAreas(config, config.next)
}

func commandMapb(config *config, args ...string) error {
	if config.previous == "" {
		return fmt.Errorf("you're on the first page")
	}

	return printLocationAreas(config, config.previous)
}

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
		if config.pokedex == nil {
			config.pokedex = make(map[string]Pokemon)
		}
		config.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

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

	return commands
}
