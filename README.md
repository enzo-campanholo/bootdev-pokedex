# Bootdev Pokedex

A CLI Pokedex built in Go that lets you explore map areas, catch Pokemon, and view their stats using the [PokeAPI](https://pokeapi.co/).

Versão em português: [README.pt.md](README.pt.md)

## Usage

Start the REPL by running:

```bash
go run .
```

## Commands

| Command | Description |
|---------|-------------|
| `map` | Display the next 20 location areas. Successive calls paginate forward. |
| `mapb` | Display the previous 20 location areas, rolling your position back. |
| `explore <area>` | List all Pokemon found in a specific location area. |
| `catch <pokemon>` | Attempt to catch a Pokemon. Higher base experience means lower catch rate. |
| `inspect <pokemon>` | View the stats and types of a caught Pokemon. |
| `pokedex` | List all Pokemon you have caught. |
| `help` | Show available commands. |
| `exit` | Exit the Pokedex CLI. |

## Design Choices

This project follows the Boot.dev Pokedex guide, but I made a few intentional changes.

### 1) Global command state instead of a `Config` struct

The course version passes a `Config` struct to every command callback. In this implementation, shared state is kept at package level (`apiClient`, `nextLocationAreasURL`, `prevLocationAreasURL`, and `pokedex`).

I chose this because:

- It removes parameters that many commands do not use.
- It makes command signatures simpler (`func(args []string) error`).
- It keeps common state access direct and readable.

### 2) Domain-focused models instead of raw API response shapes

PokeAPI responses are nested and verbose. Instead of working with those raw shapes throughout the app, this project maps responses into smaller, task-focused models.

Examples:

- We define response structs with only the fields the CLI uses; extra JSON fields are ignored during unmarshaling.
- `GetLocationAreaPokemon` returns `[]string` (Pokemon names) instead of a full nested response.
- `Pokemon` uses custom JSON unmarshaling to flatten stats and types into REPL-friendly fields.

This keeps the REPL code easier to read and makes API interaction logic easier to maintain.
