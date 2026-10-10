package cli

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/ui"
)

func commandExit(config *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, args ...string) error {
	fmt.Println("\nWelcome to the Pokedex!\nUsage:")
	for _, cmd := range config.Commands {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println()
	return nil
}

func commandMap(config *Config, args ...string) error {
	if config.Next == nil && config.Previous != nil {
		return fmt.Errorf("you're on the last page")
	}
	return fetchAndPrintLocationAreas(config, config.Next)
}

func commandMapb(config *Config, args ...string) error {
	if config.Previous == nil {
		return fmt.Errorf("you're on the first page")
	}
	return fetchAndPrintLocationAreas(config, config.Previous)
}

func fetchAndPrintLocationAreas(config *Config, pageURL *string) error {
	spinner := ui.NewSpinner("Fetching map data...")
	spinner.Start()
	resp, err := config.Client.GetLocationAreas(pageURL)
	spinner.Stop()
	if err != nil {
		return err
	}
	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}
	config.Next = resp.Next
	config.Previous = resp.Previous
	return nil
}

func commandExplore(config *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location name")
	}
	location := args[0]
	spinner := ui.NewSpinner(fmt.Sprintf("Exploring %s...", location))
	spinner.Start()
	area, err := config.Client.GetLocationArea(location)
	spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Println(ui.ColorCyan(fmt.Sprintf("Found Pokemon in %s:", location)))
	for _, enc := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	pokemonName := strings.ToLower(args[0])
	if _, exists := config.Pokedex.Get(pokemonName); exists {
		return fmt.Errorf("you have already caught %s", pokemonName)
	}
	spinner := ui.NewSpinner(fmt.Sprintf("Locating %s...", pokemonName))
	spinner.Start()
	pokemon, err := config.Client.GetPokemon(pokemonName)
	spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	baseExp := pokemon.BaseExperience
	if baseExp <= 0 {
		baseExp = 1
	}
	res := rand.Intn(baseExp)
	if res < 40 {
		fmt.Println(ui.ColorGreen(fmt.Sprintf("%s was caught!", pokemon.Name)))
		fmt.Println("You may now inspect it with the inspect command.")
		config.Pokedex.Add(pokemon)
	} else {
		fmt.Println(ui.ColorYellow(fmt.Sprintf("%s escaped!", pokemon.Name)))
	}
	return nil
}

func commandInspect(config *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	pokemonName := strings.ToLower(args[0])
	pokemon, exists := config.Pokedex.Get(pokemonName)
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", pokemon.Name, pokemon.Height, pokemon.Weight)
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", ui.ColorType(t.Type.Name))
	}
	return nil
}

func commandPokedex(config *Config, args ...string) error {
	pokemonMap := config.Pokedex.GetAll()
	if len(pokemonMap) == 0 {
		fmt.Println("Your Pokedex is empty. Go catch some Pokemon!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range pokemonMap {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

func commandRelease(config *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	pokemonName := strings.ToLower(args[0])
	pokemon, exists := config.Pokedex.Get(pokemonName)
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	config.Pokedex.Remove(pokemonName)
	fmt.Println(ui.ColorYellow(fmt.Sprintf("Bye bye, %s! %s was released back into the wild.", pokemon.Name, pokemon.Name)))
	return nil
}

func commandSummary(config *Config, args ...string) error {
	pokemonMap := config.Pokedex.GetAll()
	total := len(pokemonMap)
	fmt.Println(ui.ColorCyan("--- Trainer Summary ---"))
	fmt.Printf("Total Pokémon caught: %d\n", total)
	if total == 0 {
		fmt.Println("Catch some Pokémon first to view your statistics!")
		return nil
	}
	var heaviest, tallest pokeapi.Pokemon
	typeCounts := make(map[string]int)
	first := true
	for _, p := range pokemonMap {
		if first {
			heaviest, tallest, first = p, p, false
		} else {
			if p.Weight > heaviest.Weight { heaviest = p }
			if p.Height > tallest.Height { tallest = p }
		}
		for _, t := range p.Types { typeCounts[t.Type.Name]++ }
	}
	fmt.Printf("Heaviest Pokémon: %s (%d)\nTallest Pokémon: %s (%d)\n", heaviest.Name, heaviest.Weight, tallest.Name, tallest.Height)
	var favType string
	maxCount := 0
	for t, count := range typeCounts {
		if count > maxCount {
			maxCount, favType = count, t
		}
	}
	if favType != "" {
		fmt.Printf("Favorite Type: %s (%d)\n", ui.ColorType(favType), maxCount)
	}
	fmt.Println("Type Breakdown:")
	for t, count := range typeCounts {
		fmt.Printf("  - %s: %d\n", ui.ColorType(t), count)
	}
	return nil
}

func commandFilter(config *Config, args ...string) error {
	if len(args) == 0 { return fmt.Errorf("you must provide an elemental type") }
	targetType := strings.ToLower(args[0])
	var matches []pokeapi.Pokemon
	for _, p := range config.Pokedex.GetAll() {
		for _, t := range p.Types {
			if strings.ToLower(t.Type.Name) == targetType {
				matches = append(matches, p)
				break
			}
		}
	}
	if len(matches) == 0 {
		fmt.Printf("No %s Pokémon found in your Pokédex.\n", ui.ColorType(targetType))
		return nil
	}
	titleType := targetType
	if len(targetType) > 0 { titleType = strings.ToUpper(targetType[:1]) + targetType[1:] }
	fmt.Printf("%s Pokémon in your Pokédex:\n", titleType)
	for _, p := range matches { fmt.Printf(" - %s\n", p.Name) }
	return nil
}

func commandSave(config *Config, args ...string) error {
	if err := config.Pokedex.Save(); err != nil { return err }
	target := config.Pokedex.Filepath()
	if target == "" { target = "in-memory (no file configured)" }
	fmt.Printf("Pokédex successfully saved to %s (%d Pokémon)\n", target, config.Pokedex.Count())
	return nil
}

func commandLoad(config *Config, args ...string) error {
	if err := config.Pokedex.Load(); err != nil { return err }
	target := config.Pokedex.Filepath()
	if target == "" { target = "in-memory (no file configured)" }
	fmt.Printf("Pokédex successfully loaded from %s (%d Pokémon)\n", target, config.Pokedex.Count())
	return nil
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit":    {Name: "exit", Description: "Exit the Pokedex", Callback: commandExit},
		"help":    {Name: "help", Description: "Displays a help message", Callback: commandHelp},
		"map":     {Name: "map", Description: "Displays the next 20 locations", Callback: commandMap},
		"mapb":    {Name: "mapb", Description: "Displays the previous 20 locations", Callback: commandMapb},
		"explore": {Name: "explore", Description: "Shows all pokemon in a given area", Callback: commandExplore},
		"catch":   {Name: "catch", Description: "Attempt to catch a pokemon", Callback: commandCatch},
		"inspect": {Name: "inspect", Description: "View details about a caught pokemon", Callback: commandInspect},
		"pokedex": {Name: "pokedex", Description: "List all caught pokemon", Callback: commandPokedex},
		"save":    {Name: "save", Description: "Save your Pokédex to disk", Callback: commandSave},
		"load":    {Name: "load", Description: "Load your Pokédex from disk", Callback: commandLoad},
		"release": {Name: "release", Description: "Release a caught pokemon", Callback: commandRelease},
		"summary": {Name: "summary", Description: "Display trainer statistics", Callback: commandSummary},
		"filter":  {Name: "filter", Description: "Filter caught Pokémon by type", Callback: commandFilter},
	}
}
