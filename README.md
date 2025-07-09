# bcncli

BConomy CLI client

`bcncli` is a command-line interface for interacting with the BConomy game API. It allows you to get information about pets, eggs, profiles, factions, market trades, leaderboards, logs, and game data directly from your terminal.

## Features
Retrieve information about
* Pets
* Eggs
* User profiles
* Faction details
* Market
* Leaderboards
* Logs
* Download game data
* Search game data

## Requirements

* Go 1.24.2 or later

## Installation

Ensure you have Go installed and your `GOPATH/bin` is in your `PATH`:

```bash
go install github.com/earentir/bcncli@latest
```

## Configuration

`bcncli` requires your BConomy API key. You can provide it in three ways:

1. **Flag**: `--apikey YOUR_API_KEY`
2. **Environment Variable**: `export BCONOMYAPI=YOUR_API_KEY`
3. **Config File**: Create a `config.json` in the current directory:

   ```json
   {
     "apikey": "YOUR_API_KEY_HERE"
   }
   ```

Note: `bcncli` uses [Viper](https://github.com/spf13/viper) for configuration management and will read `config.(json|yaml|toml|...)` if present.

## Usage

```bash
bcncli [--apikey YOUR_API_KEY] <command> [flags]
```

* Run `bcncli --help` to list all commands.
* Run `bcncli <command> --help` to see flags and options for a specific command.

### Available Commands

| Command       | 
| ------------- | 
| `pet`         | 
| `egg`         | 
| `profile`     | 
| `faction`     | 
| `market`      | 
| `leaderboard` | 
| `logs`        | 
| `gamedata`    | 
| `search`      | 

### Examples

```bash
# Show your profile
bcncli profile show --user 12345

# List your pets
bcncli pet list

# Hatch an egg by ID
bcncli egg hatch --id 67890

# View market listings in "potions" category
bcncli market list --category "potions"

# View top 10 players
bcncli leaderboard list --limit 10

# Fetch recent logs
bcncli logs list --limit 20

# Export game data to JSON
bcncli gamedata export --output game_data.json

# Search game data for quests
bcncli search quests --query "Dragon Hunt"
```

*Note: Replace command examples with actual subcommand names and flags; use `--help` to discover available options.*

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Open a pull request

## License

This project is licensed under the [GPL-2.0 License](./LICENSE).
