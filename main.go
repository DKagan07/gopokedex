package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func pokedexCursor() {
	fmt.Print("Pokedex > ")
}

func main() {
	// setup vars
	commands := commands()
	r := os.Stdin
	scanner := bufio.NewScanner(r)

	pokedexHelp()
	pokedexCursor()

	config := Config{}

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
