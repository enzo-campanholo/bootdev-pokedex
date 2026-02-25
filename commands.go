package main

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokeapi"
)

type command struct {
	name        string
	description string
	callback    func(*Config) error
}

var commands map[string]command

func init() {
	commands = map[string]command{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display the help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 Location Areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 Location Areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Display the Pokemons of a Location Area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect the stats of a Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

var errExit = errors.New("exit requested")

func commandExit(_ *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return errExit
}

func commandHelp(_ *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

const (
	locationAreaLimit     = 20
	maxLocationAreaOffset = 1093
)

func commandMap(config *Config) error {
	locationAreasResponse, err := config.pokeapiClient.GetLocationAreas(config.locationAreaOffset)
	if err != nil {
		return err
	}

	locationAreas := parseLocationAreas(locationAreasResponse)
	for _, locationAreaName := range locationAreas {
		fmt.Println(locationAreaName)
	}

	config.locationAreaOffset += locationAreaLimit
	return nil
}

func commandMapb(config *Config) error {
	if config.locationAreaOffset >= locationAreaLimit {
		config.locationAreaOffset -= locationAreaLimit
	} else {
		config.locationAreaOffset = 0
	}

	locationAreasResponse, err := config.pokeapiClient.GetLocationAreas(config.locationAreaOffset)
	if err != nil {
		return err
	}

	locationAreas := parseLocationAreas(locationAreasResponse)
	for _, locationAreaName := range locationAreas {
		fmt.Println(locationAreaName)
	}

	config.locationAreaOffset += locationAreaLimit
	return nil
}

func parseLocationAreas(locationAreasResponse pokeapi.LocationAreasResponse) []string {
	locationAreas := make([]string, 0, len(locationAreasResponse.Results))
	for _, locationArea := range locationAreasResponse.Results {
		locationAreas = append(locationAreas, locationArea.Name)
	}

	return locationAreas
}

var errUnexpectedNumArgs = errors.New("unexpected number of arguments")

func commandExplore(config *Config) error {
	if len(config.arguments) != 1 {
		return errUnexpectedNumArgs
	}

	locationAreaPokemonResponse, err := config.pokeapiClient.GetLocationAreaPokemon(config.arguments[0])
	if err != nil {
		return err
	}

	pokemonList := parseLocationAreaPokemon(locationAreaPokemonResponse)
	fmt.Printf("Exploring %s...\n", config.arguments[0])
	fmt.Println("Found Pokemon:")
	for _, pokemonName := range pokemonList {
		fmt.Printf(" - %s\n", pokemonName)
	}

	return nil
}

func parseLocationAreaPokemon(locationAreaResponse pokeapi.LocationAreaResponse) []string {
	locationAreaPokemon := make([]string, 0, len(locationAreaResponse.PokemonEncounters))
	for _, pokemonEncounter := range locationAreaResponse.PokemonEncounters {
		locationAreaPokemon = append(locationAreaPokemon, pokemonEncounter.Pokemon.Name)
	}

	return locationAreaPokemon
}

func commandCatch(config *Config) error {
	if len(config.arguments) != 1 {
		return errUnexpectedNumArgs
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", config.arguments[0])

	pokemon, err := config.pokeapiClient.GetPokemon(config.arguments[0])
	if err != nil {
		return err
	}

	// This doesn't work at all, we need to think of something different
	chance := rand.IntN(pokemon.BaseExperience)
	if chance < pokemon.BaseExperience/2 {
		fmt.Printf("%s escaped!\n", config.arguments[0])
		return nil
	}

	config.pokedex[config.arguments[0]] = pokemon
	fmt.Printf("%s was caught!\n", config.arguments[0])
	return nil
}

func commandInspect(config *Config) error {
	if len(config.arguments) != 1 {
		return errUnexpectedNumArgs
	}

	pokemon, ok := config.pokedex[config.arguments[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	stats := map[string]int{}

	for _, stat := range pokemon.Stats {
		stats[stat.Stat.Name] = stat.BaseStat
	}

	fmt.Printf("Stats:\n  -hp: %d\n  -attack: %d\n  -defense: %d\n  -special-attack: %d\n  -special-defense: %d\n  -speed: %d\n", stats["hp"], stats["attack"], stats["defense"], stats["special-attack"], stats["special-defense"], stats["speed"])

	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf("  - %s\n", typ.Type.Name)
	}

	return nil
}

func commandPokedex(config *Config) error {
	fmt.Println("Your Pokedex:")
	for key, _ := range config.pokedex {
		fmt.Printf(" - %s\n", key)
	}

	return nil
}
