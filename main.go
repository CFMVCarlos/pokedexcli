package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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
		fmt.Println("Your command was:", cleanText[0])
	}
}
