package cmd

import (
	"fmt"
	"strings"

	"github.com/dorukozerr/todo-cli/internal/config"
	"github.com/dorukozerr/todo-cli/internal/fs"
	"github.com/dorukozerr/todo-cli/internal/types"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage todo groups",
	Long: `Manage todo groups:
- group --list: Show all available groups
- group --active: Show current active group
- group --switch <name>: Switch to a different group
- group --create <name>: Create a new group
- group --delete <name>: Delete a group (moves todos to default)

Without flags, shows the current active group.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		listFlag, _ := cmd.Flags().GetBool("list")
		activeFlag, _ := cmd.Flags().GetBool("active")
		switchGroup, _ := cmd.Flags().GetString("switch")
		createGroup, _ := cmd.Flags().GetString("create")
		deleteGroup, _ := cmd.Flags().GetString("delete")

		flagCount := 0

		if listFlag {
			flagCount++
		}
		if activeFlag {
			flagCount++
		}
		if switchGroup != "" {
			flagCount++
		}
		if createGroup != "" {
			flagCount++
		}
		if deleteGroup != "" {
			flagCount++
		}

		if flagCount > 1 {
			return fmt.Errorf("%sonly one flag can be used at a time%s", config.Red, config.Reset)
		}

		c, err := fs.GetConfig()
		if err != nil {
			return fmt.Errorf("%serror loading config: %v%s", config.Red, err, config.Reset)
		}

		switch {
		case listFlag:
			return handleListGroups(c)
		case activeFlag:
			return handleActiveGroup(c)
		case switchGroup != "":
			return handleSwitchGroup(c, switchGroup)
		case createGroup != "":
			return handleCreateGroup(c, createGroup)
		case deleteGroup != "":
			return handleDeleteGroup(c, deleteGroup)
		default:
			return handleDefaultGroupDisplay(c)
		}
	},
}

func handleListGroups(c *types.Config) error {
	if len(c.Groups) == 0 {
		fmt.Printf("%sNo groups found. The default group is always available.%s\n", config.Yellow, config.Reset)
		return nil
	}

	fmt.Printf("%s%sAvailable groups (%d total):%s%s\n", config.Blue, config.Bold, len(c.Groups), config.Reset, config.Reset)
	for _, group := range c.Groups {
		todoCount := countTodosInGroup(c.Todos, group.Name)
		if group.Name == c.ActiveGroup {
			fmt.Printf("  %s%s%s (%d todos) - %sACTIVE%s\n", config.Green+config.Bold, group.Name, config.Reset, todoCount, config.Cyan+config.Bold, config.Reset)
		} else {
			fmt.Printf("  %s%s%s (%d todos)\n", config.Yellow, group.Name, config.Reset, todoCount)
		}
	}
	return nil
}

func handleActiveGroup(c *types.Config) error {
	activeGroup := c.ActiveGroup
	if activeGroup == "" {
		activeGroup = "default"
	}

	todoCount := countTodosInGroup(c.Todos, activeGroup)
	incompleteCount := countIncompleteTodosInGroup(c.Todos, activeGroup)

	fmt.Printf("%sActive group:%s %s%s%s\n", config.Cyan, config.Reset, config.Green+config.Bold, activeGroup, config.Reset)
	fmt.Printf("%sTotal todos:%s %s%d%s (%s%d incomplete%s)\n",
		config.Cyan, config.Reset,
		config.Blue, todoCount, config.Reset,
		config.Yellow, incompleteCount, config.Reset)
	return nil
}

func handleSwitchGroup(c *types.Config, groupName string) error {
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return fmt.Errorf("%sgroup name cannot be empty%s", config.Red, config.Reset)
	}

	if groupName != "default" && !groupExists(c.Groups, groupName) {
		return fmt.Errorf("%sgroup '%s' does not exist. Use 'todo group --create %s' to create it first%s", config.Red, groupName, groupName, config.Reset)
	}

	c.ActiveGroup = groupName
	if err := fs.SaveConfig(c); err != nil {
		return fmt.Errorf("%serror saving config: %v%s", config.Red, err, config.Reset)
	}

	todoCount := countTodosInGroup(c.Todos, groupName)
	incompleteCount := countIncompleteTodosInGroup(c.Todos, groupName)

	fmt.Printf("%sSwitched to group '%s%s%s'%s\n", config.Green, config.Bold, groupName, config.Green, config.Reset)
	fmt.Printf("This group has %s%d%s todos (%s%d%s incomplete)\n",
		config.Blue, todoCount, config.Reset,
		config.Yellow, incompleteCount, config.Reset)
	return nil
}

func handleCreateGroup(c *types.Config, groupName string) error {
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return fmt.Errorf("%sgroup name cannot be empty%s", config.Red, config.Reset)
	}

	if groupName == "default" {
		return fmt.Errorf("%s'default' is a reserved group name%s", config.Red, config.Reset)
	}

	if groupExists(c.Groups, groupName) {
		return fmt.Errorf("%sgroup '%s' already exists%s", config.Red, groupName, config.Reset)
	}

	newGroup := types.Group{Name: groupName}
	c.Groups = append(c.Groups, newGroup)

	if err := fs.SaveConfig(c); err != nil {
		return fmt.Errorf("%serror saving config: %v%s", config.Red, err, config.Reset)
	}

	fmt.Printf("%sCreated group '%s%s%s'%s\n", config.Green, config.Bold, groupName, config.Green, config.Reset)
	fmt.Printf("Use '%stodo group --switch %s%s' to make it active\n", config.Cyan, groupName, config.Reset)
	return nil
}

func handleDeleteGroup(c *types.Config, groupName string) error {
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return fmt.Errorf("%sgroup name cannot be empty%s", config.Red, config.Reset)
	}

	if groupName == "default" {
		return fmt.Errorf("%scannot delete the default group%s", config.Red, config.Reset)
	}

	if !groupExists(c.Groups, groupName) {
		return fmt.Errorf("%sgroup '%s' does not exist%s", config.Red, groupName, config.Reset)
	}

	todoCount := countTodosInGroup(c.Todos, groupName)

	for i := range c.Todos {
		if c.Todos[i].Group == groupName {
			c.Todos[i].Group = "default"
		}
	}

	var newGroups []types.Group
	for _, group := range c.Groups {
		if group.Name != groupName {
			newGroups = append(newGroups, group)
		}
	}
	c.Groups = newGroups

	if c.ActiveGroup == groupName {
		c.ActiveGroup = "default"
		fmt.Printf("%sSwitched active group to 'default'%s\n", config.Yellow, config.Reset)
	}

	if err := fs.SaveConfig(c); err != nil {
		return fmt.Errorf("%serror saving config: %v%s", config.Red, err, config.Reset)
	}

	fmt.Printf("%sDeleted group '%s%s%s'%s\n", config.Red, config.Bold, groupName, config.Red, config.Reset)
	if todoCount > 0 {
		fmt.Printf("Moved %s%d%s todos to '%sdefault%s' group\n",
			config.Blue, todoCount, config.Reset,
			config.Green, config.Reset)
	}
	return nil
}

func handleDefaultGroupDisplay(c *types.Config) error {
	activeGroup := c.ActiveGroup
	if activeGroup == "" {
		activeGroup = "default"
	}
	fmt.Printf("%sCurrent active group:%s %s%s%s\n", config.Cyan, config.Reset, config.Green+config.Bold, activeGroup, config.Reset)

	if len(c.Groups) > 0 {
		fmt.Printf("\n%sAvailable groups:%s\n", config.Blue, config.Reset)
		for _, group := range c.Groups {
			if group.Name == c.ActiveGroup {
				fmt.Printf("  %s%s%s (%sactive%s)\n", config.Green+config.Bold, group.Name, config.Reset, config.Cyan, config.Reset)
			} else {
				fmt.Printf("  %s%s%s\n", config.Yellow, group.Name, config.Reset)
			}
		}
	}
	return nil
}

func groupExists(groups []types.Group, name string) bool {
	for _, group := range groups {
		if group.Name == name {
			return true
		}
	}
	return false
}

func countTodosInGroup(todos []types.Todo, groupName string) int {
	count := 0
	for _, todo := range todos {
		if todo.Group == groupName {
			count++
		}
	}
	return count
}

func countIncompleteTodosInGroup(todos []types.Todo, groupName string) int {
	count := 0
	for _, todo := range todos {
		if todo.Group == groupName && !todo.Completed {
			count++
		}
	}
	return count
}

func init() {
	groupCmd.Flags().BoolP("list", "l", false, "List all available groups")
	groupCmd.Flags().BoolP("active", "a", false, "Show current active group")
	groupCmd.Flags().StringP("switch", "s", "", "Switch to a different group")
	groupCmd.Flags().StringP("create", "c", "", "Create a new group")
	groupCmd.Flags().StringP("delete", "d", "", "Delete a group (moves todos to default)")
}
