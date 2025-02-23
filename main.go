package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sniidu/pokedexcli/internal/pokeapi"
	"github.com/sniidu/pokedexcli/internal/pokecache"
	"github.com/sniidu/pokedexcli/internal/shared"
)

// Commands
//
//	creating struct and declaring map for rest of functions
//	map get's filled in init
type cliCommand struct {
	name        string
	description string
	callback    func() error
	config      *shared.Config
}

var (
	commands map[string]cliCommand
	cache    = pokecache.NewCache(time.Second * 50)
)

func init() {
	mapConfig := shared.Config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous: "",
	}

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays location areas",
			callback:    commandMap,
			config:      &mapConfig,
		},
		"mapb": {
			name:        "map",
			description: "Displays location areas",
			callback:    commandMapBack,
			config:      &mapConfig,
		},
	}
}

func cleanInput(text string) []string {
	// Split text by whitespace and return lowercased in slice
	var result []string
	for _, field := range strings.Fields(text) {
		result = append(result, strings.ToLower(field))
	}
	return result
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for name, cli := range commands {
		fmt.Printf("%s: %s\n", name, cli.description)
	}
	return nil
}

func commandMap() error {
	pokeapi.Map(commands["map"].config, false, cache)
	return nil
}

func commandMapBack() error {
	pokeapi.Map(commands["map"].config, true, cache)
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())[0]
		currentCliCommand, found := commands[command]
		if !found {
			fmt.Println("Unknown command")
			continue
		}
		currentCliCommand.callback()
	}
}
