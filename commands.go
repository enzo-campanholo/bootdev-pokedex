package main

import (
	"errors"
	"fmt"

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
