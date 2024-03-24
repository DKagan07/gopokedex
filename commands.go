package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/DKagan07/gopokedex/pokeapi"
	"github.com/DKagan07/gopokedex/pokecache"
)

const pokemonApiV2 = "https://pokeapi.co/api/v2"

type Config struct {
	next    string
	prev    string
	cache   *pokecache.Cache
	pokedex map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
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
		"explore": {
			name:        "explore",
			description: "Explores the pokemon found in a particular area",
			callback:    exploreCallback,
		},
		"catch": {
			name:        "catch",
			description: "Gives you a /chance/ to add a Pokemon to your Pokedex",
			callback:    catchCallback,
		},
		"inspect": {
			name:        "inspect",
			description: "Lists some attributes of a pokemon you've caught",
			callback:    inspectCallback,
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

func helpCallback(cfg *Config, param string) error {
	if param != "" {
		return fmt.Errorf("Have a param, don't want one")
	}
	pokedexHelp()
	return nil
}

func exitCallback(cfg *Config, param string) error {
	if param != "" {
		return fmt.Errorf("Have a param, don't want one")
	}
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

func mapCallback(cfg *Config, param string) error {
	if param != "" {
		return fmt.Errorf("Have a param, don't want one")
	}

	locations := pokeapi.PokemonApiLocationResult{}

	if cfg.next == "" {
		cfg.next = fmt.Sprintf("%s/location-area", pokemonApiV2)
	}
	url := cfg.next
	v, ok := cfg.cache.Get(url)
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

func mapbCallback(cfg *Config, param string) error {
	if param != "" {
		return fmt.Errorf("Have a param, don't want one")
	}

	locations := pokeapi.PokemonApiLocationResult{}

	if cfg.prev == "" {
		return errors.New("No previous map bundle")
	}

	url := cfg.prev

	v, ok := cfg.cache.Get(url)
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

func exploreCallback(cfg *Config, param string) error {
	pokePerLocation := pokeapi.PokemonApiPokesPerLocation{}
	url := fmt.Sprintf("%s/location-area/%s", pokemonApiV2, param)

	v, ok := cfg.cache.Get(url)
	if !ok {
		bytes, err := pokeapi.PokemonApiPokePerLocationCall(url)
		if err != nil {
			return fmt.Errorf("http request to get pokesperlocation: %v", err)
		}

		cfg.cache.Add(url, bytes)

		pokePerLocation, err = pokeapi.UnmarshalToPokePerLocationResult(bytes)
		if err != nil {
			return fmt.Errorf("unmarshaling pokeperloation w/out cache: %v", err)
		}
	} else {
		var err error
		pokePerLocation, err = pokeapi.UnmarshalToPokePerLocationResult(v)
		if err != nil {
			return fmt.Errorf("unmarshaling pokeperloation w/ cache: %v", err)
		}
	}

	fmt.Printf("Exploring %s...\n", param)
	fmt.Println("Found Pokemon: ")
	for _, v := range pokePerLocation.PokemonEncounters {
		fmt.Println("- ", v.Pokemon.Name)
	}

	return nil
}

func catchCallback(cfg *Config, param string) error {
	fmt.Printf("Throwing pokeball at %s...\n", param)
	pkmn := pokeapi.Pokemon{}
	if param == "" {
		return errors.New("Not enough paramters to call 'catch'")
	}
	url := fmt.Sprintf("%s/pokemon/%s", pokemonApiV2, param)

	v, ok := cfg.cache.Get(url)
	if !ok {
		body, err := pokeapi.PokemonApiPokeInfoByNameCall(url)
		if err != nil {
			return fmt.Errorf("sending poke http call: %v", err)
		}

		cfg.cache.Add(url, body)

		pkmn, err = pokeapi.UnmarshalToPokemonInfo(body)
		if err != nil {
			return err
		}
	} else {
		var err error
		pkmn, err = pokeapi.UnmarshalToPokemonInfo(v)
		if err != nil {
			return errors.New("unmarshaling json")
		}
	}

	// Logic to determine if the user catches the pokemon
	baseExp := pkmn.BaseExperience
	randomMulti := rand.Float32()
	catchPercentage := float32(baseExp) * randomMulti

	if catchPercentage <= 50.00 {
		fmt.Printf("Caught %s!\n", param)
		if _, ok := cfg.pokedex[param]; !ok {
			cfg.pokedex[param] = pkmn
		} else {
			fmt.Printf("%s is already in your pokedex!\n", param)
		}
	} else {
		fmt.Printf("%s escaped!\n", param)
	}

	return nil
}

func inspectCallback(cfg *Config, param string) error {
	dex := cfg.pokedex
	pokeStats, ok := dex[param]
	if !ok {
		return fmt.Errorf("You have not caught %s.", param)
	}

	fmt.Println("Name:", param)
	fmt.Println("Height:", pokeStats.Height)
	fmt.Println("Weight:", pokeStats.Weight)
	fmt.Println("Stats: ")
	for _, v := range pokeStats.Stats {
		fmt.Printf("  -%s: %d\n", v.Stat.Name, v.BaseStat)
	}
	fmt.Println("Types:")
	for _, v := range pokeStats.Types {
		fmt.Println("  -", v.Type.Name)
	}

	return nil
}
