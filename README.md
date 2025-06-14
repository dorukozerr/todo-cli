# Todo CLI

Another cli todo app

## Installation

```bash
git clone git@github.com:dorukozerr/todo-cli.git ~/todo-cli &&
cd ~/todo-cli &&
go build -o todo &&
sudo mkdir -p /usr/local/bin &&
sudo mv todo /usr/local/bin &&
sudo mkdir -p /usr/local/share/zsh/site-functions &&
sudo mkdir -p /etc/bash_completion.d &&
todo completion zsh | sudo tee /usr/local/share/zsh/site-functions/_todo > /dev/null &&
todo completion bash | sudo tee /etc/bash_completion.d/todo > /dev/null &&
cd &&
rm -rf ~/todo-cli
```

## Deletion

```bash
sudo rm /usr/local/bin/todo &&
sudo rm /etc/bash_completion.d/todo &&
sudo rm /usr/local/share/zsh/site-functions/_todo
```

## **Note**: After installation, you need to restart your terminal

### Prerequisites

- Go

### Usage

App built with cobra and command descriptions are added, just use --help flag with any argument you want, `todo --help`, `todo group --help`

Example

```bash
$ todo --help
A command-line todo application with group management and priority levels

Usage:
  todo [flags]
  todo [command]

Available Commands:
  add         Add a new todo
  complete    Mark todo as completed
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a todo
  group       Manage todo groups
  help        Help about any command
  incomplete  Mark todo as incomplete
  list        List todos with filtering options
  update      Update a todo

Flags:
  -h, --help   help for todo

Use "todo [command] --help" for more information about a command.

$ todo group --help
Manage todo groups:
- group --list: Show all available groups
- group --active: Show current active group
- group --switch <name>: Switch to a different group
- group --create <name>: Create a new group
- group --delete <name>: Delete a group (moves todos to default)

Without flags, shows the current active group.

Usage:
  todo group [flags]

Flags:
  -a, --active          Show current active group
  -c, --create string   Create a new group
  -d, --delete string   Delete a group (moves todos to default)
  -h, --help            help for group
  -l, --list            List all available groups
  -s, --switch string   Switch to a different group
```
