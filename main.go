package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type command struct {
	name        string
	description string
	callback    func() error
}

var commands = map[string]command{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Display help message",
		callback:    commandHelp,
	},
}

func main() {
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

		cmd, ok := commands[userInput[0]]
		if !ok {
			fmt.Printf("Unknown command: %q\n", userInput[0])
			continue
		}

		if err := cmd.callback(); err != nil {
			if errors.Is(err, errExit) {
				return
			}
			fmt.Fprintf(os.Stderr, "command %q failed: %v\n", cmd.name, err)
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

var errExit = errors.New("exit requested")

func commandExit() error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	return errExit
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	fmt.Printf("help: Display help message\n")
	fmt.Printf("exit: Exit the Pokedex\n")
	return nil
}
