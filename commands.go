package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/DKagan07/gopokedex/pokeapi"
	"github.com/DKagan07/gopokedex/pokecache"
)

const pokemonApiV2 = "https://pokeapi.co/api/v2"

type Config struct {
	next  string
	prev  string
	cache *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

// setting up the commands and callback functions
func commands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpCallback,
		},
		"exit": {
			name:        "exit",
			description: "Exits the pokedex",
			callback:    exitCallback,
		},
		"map": {
			name:        "map",
			description: "Displays (next) 20 locations",
			callback:    mapCallback,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations",
			callback:    mapbCallback,
		},
	}
}

func pokedexHelp() {
	fmt.Println("Pokedex > These are the available commands:")
	for k, v := range commands() {
		fmt.Printf("Command: %s\n", k)
		fmt.Printf("Description: %s\n", v.description)
	}
	fmt.Println()
}

func helpCallback(cfg *Config) error {
	pokedexHelp()
	return nil
}

func exitCallback(cfg *Config) error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

func mapCallback(cfg *Config) error {
	locations := pokeapi.PokemonApiLocationResult{}

	if cfg.next == "" {
		cfg.next = fmt.Sprintf("%s/location", pokemonApiV2)
	}
	url := cfg.next
	v, ok := cfg.cache.Get(url)
	// If it's not in cache
	if !ok {
		body, err := pokeapi.PokemonApiLocationCall(url)
		if err != nil {
			return fmt.Errorf("Error with pokemon API location call: %v", err)
		}

		cfg.cache.Add(url, body)
		locations, err = pokeapi.UnmarshalToLocationResult(body)
		if err != nil {
			return fmt.Errorf("Error with unmarshalling location result: %v", err)
		}
	} else {
		fmt.Println("Used cache!")
		var err error
		locations, err = pokeapi.UnmarshalToLocationResult(v)
		if err != nil {
			return fmt.Errorf("Error with unmarshalling location result: %v", err)
		}
	}

	for _, v := range locations.Results {
		fmt.Println(v.Name)
	}

	cfg.next = locations.Next
	if locations.Previous == nil {
		cfg.prev = ""
	} else {
		cfg.prev = *locations.Previous
	}

	return nil
}

func mapbCallback(cfg *Config) error {
	locations := pokeapi.PokemonApiLocationResult{}

	if cfg.prev == "" {
		return errors.New("No previous map bundle")
	}

	url := cfg.prev

	v, ok := cfg.cache.Get(url)
	// not in get
	if !ok {
		body, err := pokeapi.PokemonApiLocationCall(url)
		if err != nil {
			return fmt.Errorf("Error with pokemon API location call: %v", err)
		}

		cfg.cache.Add(url, body)
		locations, err = pokeapi.UnmarshalToLocationResult(body)
		if err != nil {
			return fmt.Errorf("Error with unmarshalling location result: %v", err)
		}
	} else { // in get, we have []byte
		fmt.Println("used cache!")
		var err error
		locations, err = pokeapi.UnmarshalToLocationResult(v)
		if err != nil {
			return fmt.Errorf("Error with unmarshalling location result: %v", err)
		}
	}

	for _, v := range locations.Results {
		fmt.Println(v.Name)
	}

	if locations.Previous == nil {
		cfg.prev = ""
	} else {
		cfg.prev = *locations.Previous
	}

	cfg.next = locations.Next

	return nil
}
