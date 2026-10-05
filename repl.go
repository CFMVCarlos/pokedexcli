package main

import (
	"strings"
)

func cleanInput(text string) []string {
	cleanText := strings.TrimSpace(text)
	words := strings.Fields(strings.ToLower(cleanText))
	return words
}
