package cmd

import (
	"fmt"
	"strings"

	"github.com/dorukozerr/todo-cli/internal/config"
	"github.com/dorukozerr/todo-cli/internal/fs"
	"github.com/dorukozerr/todo-cli/internal/types"
	"github.com/dorukozerr/todo-cli/internal/utils"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task]",
	Short: "Add a new todo",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := args[0]

		urgency, err := cmd.Flags().GetInt("urgency")
		if err != nil {
			return
		}

		if urgency < 1 || urgency > 5 {
			fmt.Printf("%sUrgency must be between 1 and 5%s\n", config.Red, config.Reset)
			return
		}

		group, err := cmd.Flags().GetString("group")
		if err != nil {
			return
		}

		c, err := fs.GetConfig()
		if err != nil {
			fmt.Printf("%sError loading config: %v%s\n", config.Red, err, config.Reset)
			return
		}

		if group == "" {
			group = c.ActiveGroup
		}

		if group != "" {
			groupExists := false
			for _, g := range c.Groups {
				if g.Name == group {
					groupExists = true
					break
				}
			}
			if !groupExists {
				fmt.Printf("%sGroup '%s' does not exist%s\n", config.Red, group, config.Reset)
				return
			}
		}

		id := utils.GenerateNextTodoID(*c)
		newTodo := types.Todo{
			ID:        id,
			Task:      task,
			Urgency:   urgency,
			Group:     group,
			Completed: false,
		}

		c.Todos = append(c.Todos, newTodo)
		if err = fs.SaveConfig(c); err != nil {
			fmt.Printf("%sError saving config: %v%s\n", config.Red, err, config.Reset)
			return
		}

		urgencyText, urgencyColor := utils.GetUrgencyDisplay(urgency)
		groupDisplay := group
		if groupDisplay == "" {
			groupDisplay = "default"
		}

		fmt.Printf("%sAdded todo [%s%s%s]: %s%s%s\n", config.Green,
			config.Purple, id, config.Green,
			config.Bold, task, config.Reset)
		fmt.Printf("  %sUrgency:%s %s%s%s | %sGroup:%s %s%s%s\n",
			config.Cyan, config.Reset,
			urgencyColor, urgencyText, config.Reset,
			config.Cyan, config.Reset,
			config.Yellow, groupDisplay, config.Reset)
	},
}

var completeCmd = &cobra.Command{
	Use:   "complete [todo-id]",
	Short: "Mark todo as completed",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		c, err := fs.GetConfig()
		if err != nil {
			fmt.Printf("%sError loading config: %v%s\n", config.Red, err, config.Reset)
			return
		}

		for i, todo := range c.Todos {
			if todo.ID == id {
				c.Todos[i].Completed = true
				if err = fs.SaveConfig(c); err != nil {
					fmt.Printf("%sError saving config: %v%s\n", config.Red, err, config.Reset)
					return
				}
				fmt.Printf("%sCompleted todo [%s%s%s]: %s%s%s\n", config.Green,
					config.Purple, id, config.Green,
					config.Bold, todo.Task, config.Reset)
				return
			}
		}

		fmt.Printf("%sTodo with ID '%s' not found%s\n", config.Red, id, config.Reset)
	},
}

var incompleteCmd = &cobra.Command{
	Use:   "incomplete [todo-id]",
	Short: "Mark todo as incomplete",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		c, err := fs.GetConfig()
		if err != nil {
			fmt.Printf("%sError loading config: %v%s\n", config.Red, err, config.Reset)
			return
		}

		for i, todo := range c.Todos {
			if todo.ID == id {
				c.Todos[i].Completed = false
				if err = fs.SaveConfig(c); err != nil {
					fmt.Printf("%sError saving config: %v%s\n", config.Red, err, config.Reset)
					return
				}
				fmt.Printf("%sMarked todo [%s%s%s] as incomplete: %s%s%s\n", config.Yellow,
					config.Purple, id, config.Yellow,
					config.Bold, todo.Task, config.Reset)
				return
			}
		}

		fmt.Printf("%sTodo with ID '%s' not found%s\n", config.Red, id, config.Reset)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [todo-id]",
	Short: "Update a todo",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		task, _ := cmd.Flags().GetString("task")
		urgency, _ := cmd.Flags().GetInt("urgency")
		group, _ := cmd.Flags().GetString("group")
		urgencyChanged := cmd.Flags().Changed("urgency")

		if urgencyChanged && (urgency < 1 || urgency > 5) {
			fmt.Printf("%sUrgency must be between 1 and 5%s\n", config.Red, config.Reset)
			return
		}

		c, err := fs.GetConfig()
		if err != nil {
			fmt.Printf("%sError loading config: %v%s\n", config.Red, err, config.Reset)
			return
		}

		for i, todo := range c.Todos {
			if todo.ID == id {
				if group != "" {
					groupExists := false
					for _, g := range c.Groups {
						if g.Name == group {
							groupExists = true
							break
						}
					}
					if !groupExists {
						fmt.Printf("%sGroup '%s' not found%s\n", config.Red, group, config.Reset)
						return
					}
				}

				var updates []string
				if task != "" {
					c.Todos[i].Task = task
					updates = append(updates, fmt.Sprintf("task: %s%s%s", config.Bold, task, config.Reset))
				}
				if urgencyChanged {
					c.Todos[i].Urgency = urgency
					urgencyText, urgencyColor := utils.GetUrgencyDisplay(urgency)
					updates = append(updates, fmt.Sprintf("urgency: %s%s%s", urgencyColor, urgencyText, config.Reset))
				}
				if group != "" {
					c.Todos[i].Group = group
					updates = append(updates, fmt.Sprintf("group: %s%s%s", config.Yellow, group, config.Reset))
				}

				if err = fs.SaveConfig(c); err != nil {
					fmt.Printf("%sError saving config: %v%s\n", config.Red, err, config.Reset)
					return
				}

				fmt.Printf("%sUpdated todo [%s%s%s]%s\n", config.Green, config.Purple, id, config.Green, config.Reset)
				if len(updates) > 0 {
					fmt.Printf("  %sChanges:%s %s\n", config.Cyan, config.Reset, strings.Join(updates, ", "))
				}
				return
			}
		}

		fmt.Printf("%sTodo with ID '%s' not found%s\n", config.Red, id, config.Reset)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [todo-id]",
	Short: "Delete a todo",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		c, err := fs.GetConfig()
		if err != nil {
			fmt.Printf("%sError loading config: %v%s\n", config.Red, err, config.Reset)
			return
		}

		var newTodos []types.Todo
		var deletedTodo *types.Todo

		for _, todo := range c.Todos {
			if todo.ID != id {
				newTodos = append(newTodos, todo)
			} else {
				deletedTodo = &todo
			}
		}

		if deletedTodo == nil {
			fmt.Printf("%sTodo with ID '%s' not found%s\n", config.Red, id, config.Reset)
			return
		}

		c.Todos = newTodos
		if err = fs.SaveConfig(c); err != nil {
			fmt.Printf("%sError saving config: %v%s\n", config.Red, err, config.Reset)
			return
		}

		fmt.Printf("%sDeleted todo [%s%s%s]: %s%s%s\n", config.Red,
			config.Purple, id, config.Red,
			config.Bold, deletedTodo.Task, config.Reset)
	},
}

func init() {
	addCmd.Flags().IntP("urgency", "u", 1, "Set urgency level (1-5)")
	addCmd.Flags().StringP("group", "g", "", "Assign to group")

	updateCmd.Flags().StringP("task", "t", "", "Update todo task")
	updateCmd.Flags().IntP("urgency", "u", 0, "Update urgency level (1-5)")
	updateCmd.Flags().StringP("group", "g", "", "Update group assignment")
}
