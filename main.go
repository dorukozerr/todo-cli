package main

import (
	//	"flag"
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

func GetTodosPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "devtodo-cli")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}

	result := filepath.Join(configDir, "todos.json")

	return result, nil
}

func main() {
	todosPath, err := GetTodosPath()
	if err != nil {
		fmt.Printf("failed to create config directory: %v", err)
		return
	}

	//	showAll := flag.Bool("sa", false, "Show all todos")
	//	addTodo := flag.Bool("at", false, "Add a new todo")
	//
	//	flag.Parse()

	fmt.Printf("todosPath => %v\n", todosPath)
	//  fmt.Printf("*showAll => %v\n", *showAll)
	//  fmt.Printf("naddTodo => %v\n", *addTodo)
}
