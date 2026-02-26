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

var (
	apiClient            *pokeapi.Client
	nextLocationAreasURL string
	prevLocationAreasURL string
	pokedex              map[string]pokeapi.Pokemon
)

func main() {
	cache := pokecache.NewCache(5 * time.Minute)
	apiClient = pokeapi.NewClient(cache)
	pokedex = map[string]pokeapi.Pokemon{}

	scanner := bufio.NewScanner(os.Stdin)

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

		cmdName := userInput[0]
		cmd, ok := commands[cmdName]
		if !ok {
			fmt.Printf("Unknown command: %q\n", cmdName)
			continue
		}

		if err := cmd.callback(userInput[1:]); err != nil {
			if errors.Is(err, errExit) {
				break
			}
			fmt.Fprintf(os.Stderr, "command %q failed: %v\n", cmdName, err)
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
