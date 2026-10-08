package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"pokedexcli/internal/pokecache"
)

func TestMain(m *testing.M) {
	SetColorEnabled(false)
	os.Exit(m.Run())
}

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
	commands := GetCommands()

	if commands == nil {
		t.Fatal("GetCommands() returned nil")
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
			key:          "map",
			expectedName: "map",
			expectedDesc: "Displays the next 20 locations",
			expectFound:  true,
		},
		{
			key:          "mapb",
			expectedName: "mapb",
			expectedDesc: "Displays the previous 20 locations",
			expectFound:  true,
		},
		{
			key:          "explore",
			expectedName: "explore",
			expectedDesc: "Shows all pokemon in a given area",
			expectFound:  true,
		},
		{
			key:          "catch",
			expectedName: "catch",
			expectedDesc: "Attempt to catch a pokemon",
			expectFound:  true,
		},
		{
			key:          "inspect",
			expectedName: "inspect",
			expectedDesc: "View details about a caught pokemon",
			expectFound:  true,
		},
		{
			key:          "pokedex",
			expectedName: "pokedex",
			expectedDesc: "List all caught pokemon",
			expectFound:  true,
		},
		{
			key:          "save",
			expectedName: "save",
			expectedDesc: "Save your Pokédex to disk",
			expectFound:  true,
		},
		{
			key:          "load",
			expectedName: "load",
			expectedDesc: "Load your Pokédex from disk",
			expectFound:  true,
		},
		{
			key:          "release",
			expectedName: "release",
			expectedDesc: "Release a caught pokemon back into the wild",
			expectFound:  true,
		},
		{
			key:          "summary",
			expectedName: "summary",
			expectedDesc: "Display trainer statistics and collection overview",
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

	if len(commands) != 12 {
		t.Errorf("expected exactly 12 commands registered, got %d", len(commands))
	}
}

func TestCommandHelp(t *testing.T) {
	cfg := &config{
		commands: GetCommands(),
	}

	testCases := []struct {
		name string
		call func() error
	}{
		{
			name: "direct commandHelp call",
			call: func() error {
				return commandHelp(cfg)
			},
		},
		{
			name: "callback from getCommands",
			call: func() error {
				cmd, ok := GetCommands()["help"]
				if !ok {
					t.Fatal("help command not found in getCommands")
				}
				return cmd.callback(cfg)
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
	cfg := &config{
		commands: GetCommands(),
	}

	// When executed in subprocess mode, trigger the specified exit scenario and return
	mode := os.Getenv("TEST_COMMAND_EXIT_MODE")
	if mode != "" {
		switch mode {
		case "direct":
			_ = commandExit(cfg)
		case "callback":
			cmd := GetCommands()["exit"]
			_ = cmd.callback(cfg)
		case "scanner":
			scanner := bufio.NewScanner(strings.NewReader("exit\n"))
			if scanner.Scan() {
				words := cleanInput(scanner.Text())
				if len(words) > 0 {
					cmd := GetCommands()[words[0]]
					_ = cmd.callback(cfg)
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

func TestCommandMapAndMapb(t *testing.T) {
	t.Run("map on last page", func(t *testing.T) {
		cfg := &config{
			next: "",
		}
		err := commandMap(cfg)
		if err == nil || err.Error() != "you're on the last page" {
			t.Fatalf("expected 'you're on the last page', got: %v", err)
		}
	})

	t.Run("mapb on first page", func(t *testing.T) {
		cfg := &config{
			previous: "",
		}
		err := commandMapb(cfg)
		if err == nil || err.Error() != "you're on the first page" {
			t.Fatalf("expected 'you're on the first page', got: %v", err)
		}
	})

	t.Run("map and mapb successful pagination", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if strings.Contains(path, "page2") {
				resp := locationAreaResponse{
					Count: 40,
					Next:  nil,
					Results: []struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{
						{Name: "area-3", URL: "url-3"},
						{Name: "area-4", URL: "url-4"},
					},
				}
				prev := "http://" + r.Host + "/page1"
				resp.Previous = &prev
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(resp)
			} else {
				next := "http://" + r.Host + "/page2"
				resp := locationAreaResponse{
					Count: 40,
					Next:  &next,
					Results: []struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{
						{Name: "area-1", URL: "url-1"},
						{Name: "area-2", URL: "url-2"},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(resp)
			}
		}))
		defer server.Close()

		cfg := &config{
			next:     server.URL + "/page1",
			previous: "",
		}

		// First map call
		err := commandMap(cfg)
		if err != nil {
			t.Fatalf("unexpected error on commandMap: %v", err)
		}
		if cfg.next != server.URL+"/page2" {
			t.Errorf("expected cfg.next to be page2 URL, got %q", cfg.next)
		}
		if cfg.previous != "" {
			t.Errorf("expected cfg.previous to be empty, got %q", cfg.previous)
		}

		// Second map call (advance to page 2)
		err = commandMap(cfg)
		if err != nil {
			t.Fatalf("unexpected error on second commandMap: %v", err)
		}
		if cfg.next != "" {
			t.Errorf("expected cfg.next to be empty on last page, got %q", cfg.next)
		}
		if cfg.previous != server.URL+"/page1" {
			t.Errorf("expected cfg.previous to be page1 URL, got %q", cfg.previous)
		}

		// mapb call (go back to page 1)
		err = commandMapb(cfg)
		if err != nil {
			t.Fatalf("unexpected error on commandMapb: %v", err)
		}
		if cfg.next != server.URL+"/page2" {
			t.Errorf("expected cfg.next to be page2 URL after mapb, got %q", cfg.next)
		}
		if cfg.previous != "" {
			t.Errorf("expected cfg.previous to be empty after mapb, got %q", cfg.previous)
		}

		// mapb call again should fail because we are on the first page
		err = commandMapb(cfg)
		if err == nil || err.Error() != "you're on the first page" {
			t.Fatalf("expected 'you're on the first page', got: %v", err)
		}
	})
}

func TestCommandExplore(t *testing.T) {
	t.Run("missing location argument", func(t *testing.T) {
		cfg := &config{
			cache: pokecache.NewCache(5 * time.Second),
		}
		err := commandExplore(cfg)
		if err == nil {
			t.Fatal("expected error when no location argument is provided, got nil")
		}
	})

	t.Run("explore reads from cache and displays pokemon", func(t *testing.T) {
		cfg := &config{
			cache: pokecache.NewCache(5 * time.Second),
		}

		mockArea := LocationArea{
			Name: "test-area",
			PokemonEncounters: []struct {
				Pokemon struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"pokemon"`
				VersionDetails []struct {
					Version struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"version"`
					MaxChance        int `json:"max_chance"`
					EncounterDetails []struct {
						MinLevel int `json:"min_level"`
						MaxLevel int `json:"max_level"`
						Chance   int `json:"chance"`
						Method   struct {
							Name string `json:"name"`
							URL  string `json:"url"`
						} `json:"method"`
						ConditionValues []struct {
							Name string `json:"name"`
							URL  string `json:"url"`
						} `json:"condition_values"`
						PokemonDetails *string `json:"pokemon_details"`
					} `json:"encounter_details"`
				} `json:"version_details"`
			}{
				{
					Pokemon: struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{
						Name: "pikachu",
					},
				},
				{
					Pokemon: struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					}{
						Name: "bulbasaur",
					},
				},
			},
		}

		data, err := json.Marshal(mockArea)
		if err != nil {
			t.Fatalf("failed to marshal mock area: %v", err)
		}

		url := "https://pokeapi.co/api/v2/location-area/test-area"
		cfg.cache.Add(url, data)

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandExplore(cfg, "test-area")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("commandExplore returned error: %v", callErr)
		}

		output := buf.String()
		expectedSubstrings := []string{
			"Exploring test-area...",
			"Found Pokemon:",
			"- pikachu",
			"- bulbasaur",
		}
		for _, sub := range expectedSubstrings {
			if !strings.Contains(output, sub) {
				t.Errorf("expected output to contain %q, got: %q", sub, output)
			}
		}
	})
}

func TestCommandCatch(t *testing.T) {
	t.Run("missing pokemon name argument", func(t *testing.T) {
		cfg := &config{
			cache: pokecache.NewCache(5 * time.Second),
		}
		err := commandCatch(cfg)
		if err == nil {
			t.Fatal("expected error when no pokemon argument is provided, got nil")
		}
	})

	t.Run("catch with mocked cached pokemon", func(t *testing.T) {
		cfg := &config{
			cache:   pokecache.NewCache(5 * time.Second),
			pokedex: make(map[string]Pokemon),
		}

		mockPokemon := Pokemon{
			ID:             25,
			Name:           "pikachu",
			BaseExperience: 50,
			Height:         4,
			Weight:         60,
		}

		data, err := json.Marshal(mockPokemon)
		if err != nil {
			t.Fatalf("failed to marshal mock pokemon: %v", err)
		}

		url := "https://pokeapi.co/api/v2/pokemon/pikachu"
		cfg.cache.Add(url, data)

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandCatch(cfg, "pikachu")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("commandCatch returned unexpected error: %v", callErr)
		}

		output := buf.String()
		expectedPrefix := "Throwing a Pokeball at pikachu..."
		if !strings.Contains(output, expectedPrefix) {
			t.Errorf("expected output to contain %q, got: %q", expectedPrefix, output)
		}

		if strings.Contains(output, "pikachu was caught!") {
			if _, exists := cfg.pokedex["pikachu"]; !exists {
				t.Errorf("pikachu was reported caught but not found in pokedex")
			}
		} else if !strings.Contains(output, "pikachu escaped!") {
			t.Errorf("output must contain either caught or escaped message, got: %q", output)
		}
	})
}

func TestCommandInspect(t *testing.T) {
	t.Run("missing pokemon name argument", func(t *testing.T) {
		cfg := &config{
			pokedex: make(map[string]Pokemon),
		}
		err := commandInspect(cfg)
		if err == nil {
			t.Fatal("expected error when no pokemon argument is provided, got nil")
		}
	})

	t.Run("uncaught pokemon", func(t *testing.T) {
		cfg := &config{
			pokedex: make(map[string]Pokemon),
		}

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandInspect(cfg, "pidgey")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("unexpected error: %v", callErr)
		}

		output := buf.String()
		expected := "you have not caught that pokemon"
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got: %q", expected, output)
		}
	})

	t.Run("caught pokemon displays details", func(t *testing.T) {
		cfg := &config{
			pokedex: make(map[string]Pokemon),
		}

		pokemon := Pokemon{
			ID:     16,
			Name:   "pidgey",
			Height: 3,
			Weight: 18,
		}
		pokemon.Stats = []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{
				BaseStat: 40,
				Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "hp"},
			},
			{
				BaseStat: 45,
				Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "attack"},
			},
		}
		pokemon.Types = []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		}{
			{
				Slot: 1,
				Type: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "normal"},
			},
			{
				Slot: 2,
				Type: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{Name: "flying"},
			},
		}

		cfg.pokedex["pidgey"] = pokemon

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandInspect(cfg, "pidgey")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("unexpected error: %v", callErr)
		}

		output := buf.String()
		expectedStrings := []string{
			"Name: pidgey",
			"Height: 3",
			"Weight: 18",
			"Stats:",
			"  -hp: 40",
			"  -attack: 45",
			"Types:",
			"  - normal",
			"  - flying",
		}
		for _, exp := range expectedStrings {
			if !strings.Contains(output, exp) {
				t.Errorf("expected output to contain %q, but got: %q", exp, output)
			}
		}
	})
}

func TestCommandPokedex(t *testing.T) {
	t.Run("empty pokedex", func(t *testing.T) {
		cfg := &config{
			pokedex: make(map[string]Pokemon),
		}

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandPokedex(cfg)

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("unexpected error: %v", callErr)
		}

		output := buf.String()
		if !strings.Contains(output, "Your Pokedex:") {
			t.Errorf("expected output to contain 'Your Pokedex:', got: %q", output)
		}
	})

	t.Run("pokedex with caught pokemon", func(t *testing.T) {
		cfg := &config{
			pokedex: map[string]Pokemon{
				"pidgey":   {Name: "pidgey"},
				"caterpie": {Name: "caterpie"},
			},
		}

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandPokedex(cfg)

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("unexpected error: %v", callErr)
		}

		output := buf.String()
		expectedStrings := []string{
			"Your Pokedex:",
			" - pidgey",
			" - caterpie",
		}
		for _, exp := range expectedStrings {
			if !strings.Contains(output, exp) {
				t.Errorf("expected output to contain %q, got: %q", exp, output)
			}
		}
	})
}

func TestSaveAndLoadPokedex(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "pokedex-test-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	cfg := &config{
		saveFile: tmpPath,
		pokedex: map[string]Pokemon{
			"pikachu": {ID: 25, Name: "pikachu", BaseExperience: 112},
		},
	}

	// Test commandSave
	if err := commandSave(cfg); err != nil {
		t.Fatalf("commandSave returned error: %v", err)
	}

	// Create new config and load
	newCfg := &config{
		saveFile: tmpPath,
		pokedex:  make(map[string]Pokemon),
	}
	if err := commandLoad(newCfg); err != nil {
		t.Fatalf("commandLoad returned error: %v", err)
	}

	if p, ok := newCfg.pokedex["pikachu"]; !ok || p.Name != "pikachu" {
		t.Errorf("expected pikachu in loaded pokedex, got: %+v", newCfg.pokedex)
	}

	// Test loading when file does not exist
	nonExistentCfg := &config{
		saveFile: filepathJoinNonExistent(),
		pokedex:  make(map[string]Pokemon),
	}
	if err := loadPokedex(nonExistentCfg); err != nil {
		t.Errorf("expected nil error loading non-existent file, got: %v", err)
	}
}

func filepathJoinNonExistent() string {
	return filepathSafeNonExistent()
}

func filepathSafeNonExistent() string {
	return "/tmp/non-existent-pokedex-12345.json"
}

func TestCommandRelease(t *testing.T) {
	t.Run("missing pokemon name argument", func(t *testing.T) {
		cfg := &config{
			pokedex: make(map[string]Pokemon),
		}
		err := commandRelease(cfg)
		if err == nil {
			t.Fatal("expected error when no pokemon argument is provided, got nil")
		}
	})

	t.Run("release uncaught pokemon", func(t *testing.T) {
		cfg := &config{
			pokedex: make(map[string]Pokemon),
		}

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandRelease(cfg, "charizard")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("unexpected error: %v", callErr)
		}

		output := buf.String()
		expected := "you have not caught that pokemon"
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got: %q", expected, output)
		}
	})

	t.Run("release caught pokemon successfully", func(t *testing.T) {
		cfg := &config{
			pokedex: map[string]Pokemon{
				"charizard": {Name: "charizard"},
			},
		}

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandRelease(cfg, "charizard")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("unexpected error: %v", callErr)
		}

		output := buf.String()
		expected := "Bye bye, charizard! charizard was released back into the wild."
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got: %q", expected, output)
		}

		if _, exists := cfg.pokedex["charizard"]; exists {
			t.Errorf("charizard should have been removed from pokedex")
		}
	})
}

func TestCommandSummary(t *testing.T) {
	t.Run("empty pokedex summary", func(t *testing.T) {
		cfg := &config{
			pokedex: make(map[string]Pokemon),
		}

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandSummary(cfg)

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("unexpected error: %v", callErr)
		}

		output := buf.String()
		if !strings.Contains(output, "Total Pokémon caught: 0") {
			t.Errorf("expected 0 caught count, got: %q", output)
		}
	})

	t.Run("populated pokedex summary", func(t *testing.T) {
		p1 := Pokemon{
			Name:   "pidgey",
			Height: 3,
			Weight: 18,
		}
		p1.Types = []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		}{
			{Type: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "normal"}},
		}

		p2 := Pokemon{
			Name:   "gyarados",
			Height: 65,
			Weight: 2350,
		}
		p2.Types = []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		}{
			{Type: struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{Name: "water"}},
		}

		cfg := &config{
			pokedex: map[string]Pokemon{
				"pidgey":   p1,
				"gyarados": p2,
			},
		}

		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		os.Stdout = w

		callErr := commandSummary(cfg)

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		r.Close()

		if callErr != nil {
			t.Fatalf("unexpected error: %v", callErr)
		}

		output := buf.String()
		expectedSubs := []string{
			"Total Pokémon caught: 2",
			"Heaviest Pokémon: gyarados (2350)",
			"Tallest Pokémon: gyarados (65)",
			"Type Breakdown:",
			"- normal: 1",
			"- water: 1",
		}
		for _, exp := range expectedSubs {
			if !strings.Contains(output, exp) {
				t.Errorf("expected output to contain %q, but got: %q", exp, output)
			}
		}
	})
}



