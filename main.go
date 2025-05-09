package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Todo struct {
	ID        int    `json:"id"`
	Urgency   int    `json:"urgency"`
	Type      string `json:"type"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

type TodoList struct {
	Todos []Todo `json:"todos"`
}

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "devtodo-cli")

	return configDir, nil
}

func EnsureConfigDirExists() (string, error) {
	configDir, err := GetConfigDir()

	if err != nil {
		return "", err
	}

	err = os.MkdirAll(configDir, 0755)

	if err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

func GetTodosFilePath() (string, error) {
	configDir, err := GetConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "todos.json"), nil
}

// func LoadTodos() (TodoList, error) {
// 	var todoList TodoList
//
// 	todosFilePath, err := GetTodosFilePath()
//
// 	if err != nil {
// 		return todoList, err
// 	}
//
// 	if _, err := os.Stat(todosFilePath); os.IsNotExist(err) {
// 		return TodoList{Todos: []Todo{}}, nil
// 	}
//
// 	fileData, err := os.ReadFile(todosFilePath)
//
// 	if err != nil {
// 		return todoList, fmt.Errorf("failed to read todos file: %w", err)
// 	}
//
// 	if len(fileData) == 0 {
// 		return TodoList{Todos: []Todo{}}, nil
// 	}
//
// 	err = json.Unmarshal(fileData, &todoList)
//
// 	if err != nil {
// 		return todoList, fmt.Errorf("failed to parse todos file: %w", err)
// 	}
//
// 	return todoList, nil
// }

// func SaveTodos(todoList TodoList) error {
// 	configDir, err := EnsureConfigDirExists()
//
// 	if err != nil {
// 		return err
// 	}
//
// 	todosFilePath := filepath.Join(configDir, "todos.json")
//
// 	fileData, err := json.MarshalIndent(todoList, "", "  ")
//
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal todos: %w", err)
// 	}
//
// 	err = os.WriteFile(todosFilePath, fileData, 0644)
//
// 	if err != nil {
// 		return fmt.Errorf("failed to write todos file: %w", err)
// 	}
//
// 	return nil
// }

func main() {
	_, err := EnsureConfigDirExists()

	if err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)

		return
	}

	// Define command-line flags
	showAllFlag := flag.Bool("-a", false, "Show all todos including completed ones")
	testFlag := flag.Bool("-t", false, "Test Flag")

	// Parse the flags
	flag.Parse()

	args := flag.Args()

	fmt.Printf("args => %v\n", args)

	if len(args) == 0 {
		fmt.Println("default command")

		return
	}

	fmt.Printf("Show All Flag => %v", showAllFlag)
	fmt.Printf("Test Flag => %v", testFlag)
}
