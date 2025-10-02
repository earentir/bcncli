# **bcncli**

Archiving since I am not playing the game anymore.

[![Go Version](https://img.shields.io/badge/go-1.24.2-blue?logo=go)](https://go.dev/) [![License](https://img.shields.io/badge/license-GPL--2.0-blue)](LICENSE)

A fast, ergonomic command‑line client for the **[BConomy](https://bconomy.net)** game API. `bcncli` lets you fetch information about pets, eggs, factions, market, and more — all without leaving your terminal.

---

## Table of Contents

* [Features](#features)
* [Quick Start](#quick-start)
  * [Installation](#installation)
  * [Configuration](#configuration)
* [Command Reference](#command-reference)
* [Examples](#examples)
* [Building from Source](#building-from-source)
* [Shell Completions](#shell-completions)
* [Contributing](#contributing)
* [Roadmap](#roadmap)
* [License](#license)

---

## Features

| Area             | What you can do                                     |
| ---------------- | --------------------------------------------------- |
| **Pets**         | List your pets                                      |
| **Eggs**         | List eggs and inspect their metadata                |
| **Profiles**     | View any player profile by ID                       |
| **Factions**     | Get faction stats and member lists                  |
| **Market**       | Browse live listings and filter by category         |
| **Leaderboards** | Check rankings for players or factions              |
| **Logs**         | Tail or export the global event log                 |
| **Game Data**    | Download full static datasets for offline crunching |
| **Search**       | Run full‑text searches across any game resource     |

`bcncli` is built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper), so you get familiar flag handling, config files, and shell completions out of the box.

---

##  Quick Start

###  Installation

With Go >= 1.24.2 on your PATH:

```bash
# Install the latest tagged binary
$ go install github.com/earentir/bcncli@latest
```

→ The compiled executable will land in `$(go env GOPATH)/bin`. Make sure that directory is in your `PATH`.

###  Configuration

Every request to the BConomy API needs an **API key**. `bcncli` will look for it in the following order (first win):

1. `--apikey` flag
2. `BCONOMYAPI` environment variable

> **Tip:**  you can also set the `BCONOMYAPI` environment variable in your shell profile (e.g., `.bashrc`, `.zshrc`) to avoid passing the API key every time.

```bash
export BCONOMYAPI={APIKEYHERE}
```

> **Tip:**  you can also add it to a script then make it exec and finally just source it each time for the current session:

```bash
cat <<< "export BCONOMYAPI={APIKEYHERE}" > .env
source .env
```

---

##  Command Reference

```text
Usage:  bcncli [flags] <command>
```

| Command       | Description                          |
| ------------- | ------------------------------------ |
| `pet`         | Show pet information                 |
| `egg`         | Show egg information                 |
| `profile`      | Display a player profile              |
| `faction`     | Inspect factions                     |
| `market`      | List market offers with rich filters  |
| `leaderboard` | Show top players/factions            |
| `logs`        | Stream or download server logs       |
| `gamedata`    | Export BConomy static data           |
| `search`      | Free‑text search across resources    |

Run `bcncli <command> --help` for the full tree of sub‑commands and options.

---

##  Examples

| Task                            | Command                                          |
| ------------------------------- | ------------------------------------------------ |
| Show your own profile            | `bcncli profile show --user 141964`               |
| List your pets                  | `bcncli pet owned 141964`                        |
| List eggs                       | `bcncli egg owned 141964`                        |
| View potion listings            | `bcncli market list --category "potions"`        |
| Top 10 players                  | `bcncli leaderboard list --limit 10`             |
| Last 20 logs                    | `bcncli logs list --limit 20`                    |
| Export all game data            | `bcncli gamedata export --output game_data.json` |
| Search quests for “Dragon Hunt” | `bcncli search quests --query "Dragon Hunt"`     |

---

##  Building from Source

```bash
# Clone and build
$ git clone https://github.com/earentir/bcncli.git;cd bcncli;go build
```

The resulting `bcncli` binary is completely static and has **zero runtime dependencies** – ship it anywhere.

---

##  Shell Completions

Generate shell completions and drop them somewhere on your `PATH`:

```bash
# Bash
$ bcncli completion bash > /etc/bash_completion.d/bcncli

# Zsh
$ bcncli completion zsh > "${fpath[1]}/_bcncli"
```

Fish and PowerShell completions are also supported – see `bcncli completion --help`.

---

##  Contributing

1. Fork this repo and create a feature branch (`git checkout -b feature/my-awesome-thing`)
2. Commit your changes (`git commit -am 'Add awesome thing'`)
3. Push and open a pull request – don’t forget to run `go test ./...` and `go vet ./...` first

I ❤️ issues and PRs – even small docs tweaks help.

---

##  Roadmap

* [ ] Finish all the commands
* [ ] Native JSON/CSV output selectors (`--output-format`)
* [ ] CI pipeline with automated builds for Windows/macOS/Linux

Want something else? Open an [issue](https://github.com/earentir/bcncli/issues) and let me know.

---

##  License

`bcncli` is released under the [GNU General Public License v2.0](LICENSE).
