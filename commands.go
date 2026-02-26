package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

type command struct {
	description string
	callback    func(args []string) error
}

var commands map[string]command

func init() {
	commands = map[string]command{
		"exit": {
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			description: "Display the help message",
			callback:    commandHelp,
		},
		"map": {
			description: "Display the next 20 Location Areas",
			callback:    commandMap,
		},
		"mapb": {
			description: "Display the previous 20 Location Areas",
			callback:    commandMapb,
		},
		"explore": {
			description: "Display the Pokemons of a Location Area",
			callback:    commandExplore,
		},
		"catch": {
			description: "Try to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			description: "Inspect the stats of a Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			description: "List all caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

var errExit = errors.New("exit requested")

func commandExit(_ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return errExit
}

func commandHelp(_ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for name, cmd := range commands {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}

	return nil
}

func listLocationAreas(pageURL string) error {
	areas, next, prev, err := apiClient.GetLocationAreas(pageURL)
	if err != nil {
		return err
	}

	for _, area := range areas {
		fmt.Println(area.Name)
	}

	nextLocationAreasURL = next
	prevLocationAreasURL = prev
	return nil
}

func commandMap(_ []string) error {
	return listLocationAreas(nextLocationAreasURL)
}

func commandMapb(_ []string) error {
	if prevLocationAreasURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	return listLocationAreas(prevLocationAreasURL)
}

var errUnexpectedNumArgs = errors.New("unexpected number of arguments")

func commandExplore(args []string) error {
	if len(args) != 1 {
		return errUnexpectedNumArgs
	}

	names, err := apiClient.GetLocationAreaPokemon(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, name := range names {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

func commandCatch(args []string) error {
	if len(args) != 1 {
		return errUnexpectedNumArgs
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])

	pokemon, err := apiClient.GetPokemon(args[0])
	if err != nil {
		return err
	}

	chance := rand.IntN(pokemon.BaseExperience)
	if chance < pokemon.BaseExperience/2 {
		fmt.Printf("%s escaped!\n", args[0])
		return nil
	}

	pokedex[args[0]] = pokemon
	fmt.Printf("%s was caught!\n", args[0])
	return nil
}

func commandInspect(args []string) error {
	if len(args) != 1 {
		return errUnexpectedNumArgs
	}

	pokemon, ok := pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf("  - %s\n", typ)
	}

	return nil
}

func commandPokedex(_ []string) error {
	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}
