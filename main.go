package main

import (
	// "bufio"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	Bold      = "\033[1m"
	Underline = "\033[4m"
)

func printColored(color, text string) {
	fmt.Printf("%s%s%s\n", color, text, Reset)
}

func clearLine() {
	fmt.Print("\033[2K\r")
}

func moveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

func getUrgencyWithColors() int {
	fmt.Printf("%sSelect urgency (1-5): %s", Cyan, Reset)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		urgency, err := strconv.Atoi(input)

		if err != nil || urgency < 1 || urgency > 5 {
			moveCursorUp(1)
			clearLine()
			fmt.Printf("%sInvalid! Select urgency (1-5): %s", Red, Reset)
			continue
		}

		moveCursorUp(1)
		clearLine()
		fmt.Printf("%s✓ Urgency set to: %d%s\n", Green, urgency, Reset)
		return urgency
	}
	return 0
}

type Todo struct {
	ID        int    `json:"id"`
	Group     string `json:"group"`
	Urgency   int    `json:"urgency"`
	Type      string `json:"type"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

type Config struct {
	Group string `json:"group"`
	Todos []Todo `json:"todos"`
}

func GetConfig() ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []byte{}, err
	}
	configDir := filepath.Join(homeDir, ".config", "todo-cli")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return []byte{}, err
	}
	configPath := filepath.Join(configDir, "todos.json")
	_, err = os.Stat(configPath)
	if err != nil {
		config, err := json.MarshalIndent(Config{
			Group: "",
			Todos: []Todo{},
		}, "", "  ")
		return config, err
	}
	config, err := os.ReadFile(configPath)
	return config, err
}

func main() {
	config, err := GetConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
		return
	}

	// showAll := flag.Bool("sa", false, "Show all todos")
	// group := flag.String("g", "", "Target todo group")
	// addTodo := flag.Bool("at", false, "Add a new todo")

	flag.Parse()

	var c Config
	err = json.Unmarshal(config, &c)
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v", err)
		return
	}

	if c.Group == "" {
		fmt.Println("group is empty")
	} else {
		fmt.Println("group is not empty")
	}
	urgency := getUrgencyWithColors()

	fmt.Println(urgency)
}
