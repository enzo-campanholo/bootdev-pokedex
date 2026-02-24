package main

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokeapi"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokeapiClient := pokeapi.NewClient()

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

		if err := cmd.callback(pokeapiClient); err != nil {
			if errors.Is(err, errExit) {
				break
			}
			fmt.Fprintf(os.Stderr, "command %q failed: %v\n", cmd.name, err)
		}
	}

	os.Exit(0)
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
