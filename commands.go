package main

import (
	"errors"
	"fmt"

	"github.com/enzo-campanholo/bootdev-pokedex/internal/pokeapi"
)

type command struct {
	name        string
	description string
	callback    func(*pokeapi.Client) error
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
	}
}

var errExit = errors.New("exit requested")

func commandExit(_ *pokeapi.Client) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	return errExit
}

func commandHelp(_ *pokeapi.Client) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

const (
	locationAreaLimit     = 20
	maxLocationAreaOffset = 1093
)

var locationAreaOffset int = 0

func commandMap(pokeapiClient *pokeapi.Client) error {
	locationAreasResponse, err := pokeapiClient.GetLocationAreas(locationAreaOffset)
	if err != nil {
		return err
	}

	locationAreas := parseLocationAreas(locationAreasResponse)
	for _, locationAreaName := range locationAreas {
		fmt.Printf("%s\n", locationAreaName)
	}

	locationAreaOffset += locationAreaLimit
	return nil
}

func commandMapb(pokeapiClient *pokeapi.Client) error {
	if locationAreaOffset >= locationAreaLimit {
		locationAreaOffset -= locationAreaLimit
	} else {
		locationAreaOffset = 0
	}

	locationAreasResponse, err := pokeapiClient.GetLocationAreas(locationAreaOffset)
	if err != nil {
		return err
	}

	locationAreas := parseLocationAreas(locationAreasResponse)
	for _, locationAreaName := range locationAreas {
		fmt.Printf("%s\n", locationAreaName)
	}

	locationAreaOffset += locationAreaLimit
	return nil
}

func parseLocationAreas(locationAreasResponse pokeapi.LocationAreasResponse) []string {
	locationAreas := make([]string, 0, len(locationAreasResponse.Results))
	for _, locationArea := range locationAreasResponse.Results {
		locationAreas = append(locationAreas, locationArea.Name)
	}

	return locationAreas
}
