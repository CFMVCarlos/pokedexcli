package cli

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur SQUIRTLE",
			expected: []string{"charmander", "bulbasaur", "squirtle"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("expected %v, got %v", c.expected, actual)
			continue
		}
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("expected %v, got %v", c.expected, actual)
		}
	}
}
