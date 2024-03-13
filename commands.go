package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const pokemonApiV2 = "https://pokeapi.co/api/v2"

type PokemonApiLocationResult struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	next string
	prev string
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
	if cfg.next == "" {
		cfg.next = fmt.Sprintf("%s/location", pokemonApiV2)
	}
	url := cfg.next
	locations, err := pokemonApiLocationCall(url)
	if err != nil {
		return fmt.Errorf("Error with pokemon API location call: %v", err)
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
	if cfg.prev == "" {
		return errors.New("No previous map bundle")
	}

	url := cfg.prev
	locations, err := pokemonApiLocationCall(url)
	if err != nil {
		return fmt.Errorf("Error with pokemon API location call: %v", err)
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

// Api calls
func pokemonApiLocationCall(url string) (PokemonApiLocationResult, error) {
	var locations PokemonApiLocationResult

	res, err := http.Get(url)
	if err != nil {
		return locations, fmt.Errorf("Failed HTTP request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return locations, errors.New("Error with ReadAll")
	}

	if err := json.Unmarshal(body, &locations); err != nil {
		return locations, errors.New("Error unmarshaling json")
	}

	return locations, nil
}
