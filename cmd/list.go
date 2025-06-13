package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dorukozerr/todo-cli/internal/config"
	"github.com/dorukozerr/todo-cli/internal/fs"
	"github.com/dorukozerr/todo-cli/internal/types"
	"github.com/dorukozerr/todo-cli/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List todos with filtering options",
	Long: `List todos with filtering options:
- Default: Shows incomplete todos from active group
- --all: Shows all todos from active group
- --all-groups: Shows incomplete todos from all groups
- --all --all-groups: Shows all todos from all groups`,
	Run: func(cmd *cobra.Command, args []string) {
		showAll, _ := cmd.Flags().GetBool("all")
		allGroups, _ := cmd.Flags().GetBool("all-groups")

		c, err := fs.GetConfig()
		if err != nil {
			fmt.Printf("%sError loading config: %v%s\n", config.Red, err, config.Reset)
			return
		}

		if len(c.Todos) == 0 {
			fmt.Printf("%sNo todos found%s\n", config.Yellow, config.Reset)
			return
		}

		filteredTodos := filterTodos(c.Todos, c.ActiveGroup, showAll, allGroups)

		if len(filteredTodos) == 0 {
			displayEmptyMessage(showAll, allGroups, c.ActiveGroup)
			return
		}

		sort.Slice(filteredTodos, func(i, j int) bool {
			if filteredTodos[i].Completed != filteredTodos[j].Completed {
				return !filteredTodos[i].Completed
			}
			return filteredTodos[i].Urgency > filteredTodos[j].Urgency
		})

		displayHeader(showAll, allGroups, c.ActiveGroup)

		if allGroups {
			displayTodosByGroup(filteredTodos)
		} else {
			displayTodosList(filteredTodos)
		}
	},
}

func filterTodos(todos []types.Todo, activeGroup string, showAll, allGroups bool) []types.Todo {
	var filtered []types.Todo

	for _, todo := range todos {
		if !showAll && todo.Completed {
			continue
		}
		if !allGroups && todo.Group != activeGroup {
			continue
		}
		filtered = append(filtered, todo)
	}

	return filtered
}

func displayEmptyMessage(showAll, allGroups bool, activeGroup string) {
	status := "incomplete todos"
	if showAll {
		status = "todos"
	}

	if allGroups {
		fmt.Printf("%sNo %s found in any group%s\n", config.Yellow, status, config.Reset)
	} else {
		groupName := activeGroup
		if groupName == "" {
			groupName = "default"
		}
		fmt.Printf("%sNo %s found in group '%s'%s\n", config.Yellow, status, groupName, config.Reset)
	}
}

func displayHeader(showAll, allGroups bool, activeGroup string) {
	status := "Incomplete todos"
	if showAll {
		status = "All todos"
	}

	scope := "from all groups"
	if !allGroups {
		groupName := activeGroup
		if groupName == "" {
			groupName = "default"
		}
		scope = fmt.Sprintf("from group '%s'", groupName)
	}

	fmt.Printf("\n%s%s %s:%s\n", config.Blue+config.Bold, status, scope, config.Reset)
	fmt.Println(strings.Repeat("=", 40))
}

func displayTodosByGroup(todos []types.Todo) {
	todoGroups := make(map[string][]types.Todo)
	for _, todo := range todos {
		groupName := todo.Group
		if groupName == "" {
			groupName = "default"
		}
		todoGroups[groupName] = append(todoGroups[groupName], todo)
	}

	first := true
	for groupName, groupTodos := range todoGroups {
		if !first {
			fmt.Println()
		}
		first = false

		fmt.Printf("\n%sGroup: %s%s\n", config.Cyan+config.Bold, groupName, config.Reset)
		fmt.Println(strings.Repeat("-", 20))
		displayTodosList(groupTodos)
	}
}

func displayTodosList(todos []types.Todo) {
	for _, todo := range todos {
		status := "[ ]"
		statusColor := config.Yellow
		if todo.Completed {
			status = "[x]"
			statusColor = config.Green
		}

		urgencyText, urgencyColor := utils.GetUrgencyDisplay(todo.Urgency)
		fmt.Printf("%s%s%s [%s%s%s] %s%s%s %s\n",
			statusColor, status, config.Reset,
			config.Purple, todo.ID, config.Reset,
			urgencyColor, urgencyText, config.Reset,
			todo.Task)
	}
}

func init() {
	listCmd.Flags().BoolP("all", "a", false, "Show completed and incomplete todos")
	listCmd.Flags().Bool("all-groups", false, "Show todos from all groups")
}
