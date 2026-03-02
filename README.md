# Bootdev Pokedex

A CLI Pokedex built in Go that lets you explore map areas, catch Pokemon, and view their stats using the [PokeAPI](https://pokeapi.co/).

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

