package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "mixed case with surrounding whitespace",
			input:    "  Hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "multiple spaces between words",
			input:    "  Mamma   Mia               ",
			expected: []string{"mamma", "mia"},
		},
		{
			name:     "all uppercase",
			input:    "PIKACHU CHARMANDER BULBASAUR",
			expected: []string{"pikachu", "charmander", "bulbasaur"},
		},
		{
			name:     "tabs and newlines as whitespace",
			input:    "charmander\t\nbulbasaur \t squirtle\n",
			expected: []string{"charmander", "bulbasaur", "squirtle"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only spaces",
			input:    "      ",
			expected: []string{},
		},
		{
			name:     "only whitespace characters",
			input:    " \t \n \r ",
			expected: []string{},
		},
		{
			name:     "single word with whitespace",
			input:    "   Pikachu   ",
			expected: []string{"pikachu"},
		},
		{
			name:     "punctuation and special characters preserved",
			input:    "  hello, world! 123 #pokemon  ",
			expected: []string{"hello,", "world!", "123", "#pokemon"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := cleanInput(c.input)
			if len(actual) != len(c.expected) {
				t.Fatalf("cleanInput(%q) length = %d, expected %d (got %v, expected %v)",
					c.input, len(actual), len(c.expected), actual, c.expected)
			}

			for i := range actual {
				if actual[i] != c.expected[i] {
					t.Errorf("cleanInput(%q)[%d] = %q, expected %q",
						c.input, i, actual[i], c.expected[i])
				}
			}
		})
	}
}
