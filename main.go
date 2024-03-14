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
		// have to string parse "val" here
		strs := strings.Split(val, " ")
		switch strs[0] {
		case "exit":
			if len(strs) > 1 {
				fmt.Println("Too many arguments in 'exit' command")
				break
			}
			if err := commands[val].callback(&config, ""); err != nil {
				fmt.Println("ERROR with exit: ", err)
			}
		case "help":
			if len(strs) > 1 {
				fmt.Println("Too many arguments in 'help' command")
				break
			}
			if err := commands[val].callback(&config, ""); err != nil {
				fmt.Println("ERROR with help: ", err)
			}
		case "map":
			if len(strs) > 1 {
				fmt.Println("Too many arguments in 'map' command")
				break
			}
			if err := commands[val].callback(&config, ""); err != nil {
				fmt.Println("ERROR with map: ", err)
			}
		case "mapb":
			if len(strs) > 1 {
				fmt.Println("Too many arguments in 'mapb' command")
				break
			}
			if err := commands[val].callback(&config, ""); err != nil {
				fmt.Println("ERROR with mapb: ", err)
			}
		case "explore":
			if len(strs) > 2 {
				fmt.Println("Too many arguments in 'explore' command")
				break
			}
			city := strs[1]
			if err := commands[strs[0]].callback(&config, city); err != nil {
				fmt.Println("ERROR with explore: ", err)
			}
		default:
			fmt.Println("Unknown command, see 'help' for available commands")
		}
		pokedexCursor()
	}
}
