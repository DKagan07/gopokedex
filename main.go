package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DKagan07/gopokedex/pokecache"
)

func pokedexCursor() {
	fmt.Print("Pokedex > ")
}

func main() {
	// setup vars
	commands := commands()
	r := os.Stdin
	scanner := bufio.NewScanner(r)
	dur := time.Minute

	pokedexHelp()
	pokedexCursor()

	cache := pokecache.NewCache(dur)
	config := Config{
		cache: &cache,
	}

	// the logic
	for scanner.Scan() {
		val := strings.ToLower(scanner.Text())
		switch val {
		case "exit":
			if err := commands[val].callback(&config); err != nil {
				fmt.Println("ERROR with exit: ", err)
			}
		case "help":
			if err := commands[val].callback(&config); err != nil {
				fmt.Println("ERROR with help: ", err)
			}
		case "map":
			if err := commands[val].callback(&config); err != nil {
				fmt.Println("ERROR with map: ", err)
			}
		case "mapb":
			if err := commands[val].callback(&config); err != nil {
				fmt.Println("ERROR with mapb: ", err)
			}
		default:
			fmt.Println("Unknown command, see 'help' for available commands")
		}
		pokedexCursor()
	}
}
