# Pokedex CLI

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.27+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)
![Curriculum](https://img.shields.io/badge/Boot.dev-Backend%20Track-805AD5?style=for-the-badge)
![API](https://img.shields.io/badge/PokeAPI-v2-EF5350?style=for-the-badge)
![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen?style=for-the-badge)

**An interactive command-line Pokédex REPL built in Go, powered by PokeAPI with a thread-safe, concurrent TTL cache.**

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage-and-commands) • [Architecture](#-project-architecture) • [Testing](#-testing--verification)

</div>

---

## ⚡ Features

- 🎮 **Interactive REPL Loop:** Fast, intuitive command prompt with input sanitization and command routing.
- 🗺️ **World & Area Navigation:** Paginated exploration of Pokémon location areas (`map` / `mapb`).
- 🔍 **Area Exploration:** Discover all Pokémon native to a specific location area (`explore`).
- 🎯 **Catch Mechanics:** Realistic catch probability scaled dynamically against a Pokémon's `base_experience` (`catch`).
- 📊 **Detailed Inspection:** View captured Pokémon attributes, including height, weight, base stats, and types (`inspect`).
- 📖 **Pokédex Ledger:** Real-time inventory tracking of all captured Pokémon (`pokedex`).
- ⚡ **Thread-Safe In-Memory Cache:** Custom `pokecache` package with background goroutine reaping based on configurable TTL intervals to minimize redundant network requests.

---

## 🚀 Installation

### Prerequisites

- [Go](https://go.dev/dl/) **1.22+** (configured for Go 1.27)
- `git`

### Clone and Run

```bash
# Clone the repository
git clone https://github.com/CFMVCarlos/pokedexcli.git
cd pokedexcli

# Run directly
go run .

# Or build the binary executable
go build -o pokedexcli
./pokedexcli
```

---

## 🕹️ Usage and Commands

When launched, the program starts the `Pokedex > ` interactive prompt.

| Command | Arguments | Description | Example |
| :--- | :--- | :--- | :--- |
| `help` | _None_ | Displays the command menu and descriptions | `help` |
| `exit` | _None_ | Closes the Pokédex and terminates the program | `exit` |
| `map` | _None_ | Displays the next 20 location areas from PokeAPI | `map` |
| `mapb` | _None_ | Displays the previous 20 location areas | `mapb` |
| `explore` | `<area_name>` | Lists all Pokémon available in the given location area | `explore pastoria-city-area` |
| `catch` | `<pokemon_name>` | Attempts to catch a Pokémon and adds it to your Pokédex | `catch pikachu` |
| `inspect` | `<pokemon_name>` | Displays stats, dimensions, and types of a caught Pokémon | `inspect pikachu` |
| `pokedex` | _None_ | Prints a list of all caught Pokémon in your Pokédex | `pokedex` |

### Example Session

```text
Welcome to the Pokedex!
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...
Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 - gyarados
 - remoraid

Pokedex > catch remoraid
Throwing a Pokeball at remoraid...
remoraid was caught!
You may now inspect it with the inspect command.

Pokedex > inspect remoraid
Name: remoraid
Height: 6
Weight: 120
Stats:
  -hp: 35
  -attack: 65
  -defense: 35
  -special-attack: 65
  -special-defense: 35
  -speed: 65
Types:
  - water

Pokedex > pokedex
Your Pokedex:
 - remoraid

Pokedex > exit
Closing the Pokedex... Goodbye!
```

---

## 🏗️ Project Architecture

```text
pokedexcli/
├── main.go                     # Application entrypoint & dependency setup
├── config.go                   # State types (config, cliCommand, Pokemon, LocationArea)
├── repl.go                     # REPL input scanner, dispatch loop, and command handlers
├── repl_test.go                # Comprehensive unit test suite for REPL and commands
├── internal/
│   └── pokecache/
│       ├── pokecache.go        # Thread-safe in-memory TTL cache with mutexes & ticker reap loop
│       └── pokecache_test.go   # Concurrency and expiration tests for pokecache
├── go.mod                      # Go module definition
└── README.md                   # Project documentation
```

### Technical Highlights

- **`pokecache.Cache`**: Uses `sync.Mutex` to protect the underlying map against data races. A background goroutine (`time.NewTicker`) reaps expired entries based on a configurable interval.
- **Dynamic Catch Rates**: Employs `math/rand` against Pokémon `base_experience` to balance encounter difficulty.
- **Standard Library Only**: Built cleanly using Go's standard library (`net/http`, `sync`, `time`, `encoding/json`, `bufio`, `strings`, `math/rand`).

---

## 🧪 Testing & Verification

Run the full automated test suite:

```bash
# Run all unit tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with the Go race detector enabled
go test -race ./...
```

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

## 🎓 Acknowledgements

- Built as part of the [Boot.dev](https://boot.dev) Backend Engineering curriculum.
- Pokémon data provided by [PokeAPI](https://pokeapi.co/).
