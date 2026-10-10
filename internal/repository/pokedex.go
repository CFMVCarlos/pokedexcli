package repository

import (
	"encoding/json"
	
	"os"
	"sync"

	"pokedexcli/internal/pokeapi"
)

type Pokedex struct {
	filepath string
	pokemon  map[string]pokeapi.Pokemon
	mu       sync.RWMutex
}

func NewPokedex(filepath string) *Pokedex {
	return &Pokedex{
		filepath: filepath,
		pokemon:  make(map[string]pokeapi.Pokemon),
	}
}

func (p *Pokedex) Load() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.filepath == "" {
		return nil
	}
	data, err := os.ReadFile(p.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var loaded map[string]pokeapi.Pokemon
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}
	if loaded != nil {
		p.pokemon = loaded
	} else {
		p.pokemon = make(map[string]pokeapi.Pokemon)
	}
	return nil
}

func (p *Pokedex) Save() error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.filepath == "" {
		return nil
	}
	data, err := json.MarshalIndent(p.pokemon, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.filepath, data, 0644)
}

func (p *Pokedex) Add(pokemon pokeapi.Pokemon) error {
	p.mu.Lock()
	p.pokemon[pokemon.Name] = pokemon
	p.mu.Unlock()
	return p.Save()
}

func (p *Pokedex) Remove(name string) error {
	p.mu.Lock()
	delete(p.pokemon, name)
	p.mu.Unlock()
	return p.Save()
}

func (p *Pokedex) Get(name string) (pokeapi.Pokemon, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pokemon, exists := p.pokemon[name]
	return pokemon, exists
}

func (p *Pokedex) GetAll() map[string]pokeapi.Pokemon {
	p.mu.RLock()
	defer p.mu.RUnlock()
	copy := make(map[string]pokeapi.Pokemon, len(p.pokemon))
	for k, v := range p.pokemon {
		copy[k] = v
	}
	return copy
}

func (p *Pokedex) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.pokemon)
}

func (p *Pokedex) Filepath() string {
	return p.filepath
}

// Ensure the map can be injected directly for tests
func (p *Pokedex) SetPokemonMap(pm map[string]pokeapi.Pokemon) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.pokemon = pm
}
