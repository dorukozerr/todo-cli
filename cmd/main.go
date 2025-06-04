package main

import (
	"fmt"

	"github.com/dorukozerr/todo-cli/internal/ui"
)

func main() {
	err := ui.RenderGroupsMenu()
	if err != nil {
		fmt.Printf("Error on group menu: %v", err)
	}

	// showAll := flag.Bool("sa", false, "Show all todos")
	// group := flag.String("g", "", "Target todo group")
	// addTodo := flag.Bool("at", false, "Add a new todo")

	// flag.Parse()

	// if c.Group == "" {
	// 	fmt.Println("group is empty")
	// } else {
	// 	fmt.Println("group is not empty")
	// }
}
