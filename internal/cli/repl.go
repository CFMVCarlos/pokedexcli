package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"pokedexcli/internal/ui"
)

func StartRepl(config *Config) {
	config.Pokedex.Load()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(ui.ColorCyan("Pokedex > "))
		if !scanner.Scan() { break }
		cleanText := cleanInput(scanner.Text())
		if len(cleanText) == 0 { continue }
		cmd, ok := config.Commands[cleanText[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		args := []string{}
		if len(cleanText) > 1 { args = cleanText[1:] }
		if err := cmd.Callback(config, args...); err != nil { fmt.Println(err) }
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(strings.TrimSpace(text)))
}
