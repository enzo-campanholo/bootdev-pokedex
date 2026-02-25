// Pokedex is a CLI tool for browsing Pokemon location areas using the PokeAPI.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokeapi"
	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokecache"
)

// Config holds shared state passed to every command callback.
type Config struct {
	pokeapiClient      *pokeapi.Client
	locationAreaOffset int
	arguments          []string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Minute)
	config := Config{
		pokeapiClient: pokeapi.NewClient(cache),
	}

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
			}
			return
		}

		userInput := cleanInput(scanner.Text())
		if len(userInput) == 0 {
			continue
		}

		cmd, ok := commands[userInput[0]]
		if !ok {
			fmt.Printf("Unknown command: %q\n", userInput[0])
			continue
		}

		config.arguments = userInput[1:]

		if err := cmd.callback(&config); err != nil {
			if errors.Is(err, errExit) {
				break
			}
			fmt.Fprintf(os.Stderr, "command %q failed: %v\n", cmd.name, err)
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
