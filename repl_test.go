package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
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

func TestGetCommands(t *testing.T) {
	commands := getCommands()

	if commands == nil {
		t.Fatal("getCommands() returned nil")
	}

	expectedCommands := []struct {
		key          string
		expectedName string
		expectedDesc string
		expectFound  bool
	}{
		{
			key:          "exit",
			expectedName: "exit",
			expectedDesc: "Exit the Pokedex",
			expectFound:  true,
		},
		{
			key:          "help",
			expectedName: "help",
			expectedDesc: "Displays a help message",
			expectFound:  true,
		},
		{
			key:         "invalid",
			expectFound: false,
		},
		{
			key:         "EXIT",
			expectFound: false,
		},
		{
			key:         "",
			expectFound: false,
		},
	}

	for _, tc := range expectedCommands {
		t.Run("key_"+tc.key, func(t *testing.T) {
			cmd, ok := commands[tc.key]
			if ok != tc.expectFound {
				t.Fatalf("commands[%q] found = %v, expected %v", tc.key, ok, tc.expectFound)
			}

			if tc.expectFound {
				if cmd.name != tc.expectedName {
					t.Errorf("cmd.name = %q, expected %q", cmd.name, tc.expectedName)
				}
				if cmd.description != tc.expectedDesc {
					t.Errorf("cmd.description = %q, expected %q", cmd.description, tc.expectedDesc)
				}
				if cmd.callback == nil {
					t.Errorf("cmd.callback is nil for command %q", tc.key)
				}
			}
		})
	}

	if len(commands) != 2 {
		t.Errorf("expected exactly 2 commands registered, got %d", len(commands))
	}
}

func TestCommandHelp(t *testing.T) {
	testCases := []struct {
		name string
		call func() error
	}{
		{
			name: "direct commandHelp call",
			call: commandHelp,
		},
		{
			name: "callback from getCommands",
			call: func() error {
				cmd, ok := getCommands()["help"]
				if !ok {
					t.Fatal("help command not found in getCommands")
				}
				return cmd.callback()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			os.Stdout = w

			callErr := tc.call()

			// Restore stdout and read output
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			r.Close()

			if callErr != nil {
				t.Errorf("commandHelp() returned unexpected error: %v", callErr)
			}

			output := buf.String()
			expectedSubstrings := []string{
				"exit: Exit the Pokedex",
				"help: Displays a help message",
			}

			for _, sub := range expectedSubstrings {
				if !strings.Contains(output, sub) {
					t.Errorf("expected help output to contain %q, but got: %q", sub, output)
				}
			}
		})
	}
}

func TestCommandExit(t *testing.T) {
	// When executed in subprocess mode, trigger the specified exit scenario and return
	mode := os.Getenv("TEST_COMMAND_EXIT_MODE")
	if mode != "" {
		switch mode {
		case "direct":
			_ = commandExit()
		case "callback":
			cmd := getCommands()["exit"]
			_ = cmd.callback()
		case "scanner":
			scanner := bufio.NewScanner(strings.NewReader("exit\n"))
			if scanner.Scan() {
				words := cleanInput(scanner.Text())
				if len(words) > 0 {
					cmd := getCommands()[words[0]]
					_ = cmd.callback()
				}
			}
		}
		return
	}

	cases := []struct {
		name string
		mode string
	}{
		{
			name: "direct commandExit execution",
			mode: "direct",
		},
		{
			name: "callback from getCommands exit",
			mode: "callback",
		},
		{
			name: "input scanner triggering exit command",
			mode: "scanner",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestCommandExit$")
			cmd.Env = append(os.Environ(), "TEST_COMMAND_EXIT_MODE="+tc.mode)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("expected commandExit to exit cleanly (code 0), got: %v, output: %s", err, string(out))
			}

			expectedMessage := "Closing the Pokedex... Goodbye!"
			if !strings.Contains(string(out), expectedMessage) {
				t.Errorf("expected output to contain %q, got: %q", expectedMessage, string(out))
			}
		})
	}
}
