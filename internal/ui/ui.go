package ui

import (
	"bufio"
	"fmt"
	"os"
	//"os/exec"
	"strings"

	"github.com/dorukozerr/todo-cli/internal/config"
	"github.com/dorukozerr/todo-cli/internal/fs"
	"github.com/dorukozerr/todo-cli/internal/utils"
)

func clearLines() {
	fmt.Print("\033[0J")
}

func moveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

func AddNewGroup() error {
	config, err := fs.GetConfig()
	if err != nil {
		return fmt.Errorf("Could not read config file")
	}

	fmt.Printf("%sCreate a new todo group%s\n%sEnter new group name: %s", colors.Yellow, colors.Reset, colors.Green, colors.Reset)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		groupName := strings.TrimSpace(scanner.Text())

		if len(groupName) < 3 {
			moveCursorUp(2)
			clearLines()
			fmt.Printf("%sError: Group name can be minumum 3 characters!%s\n%sEnter new group name: %s", colors.Yellow, colors.Reset, colors.Green, colors.Reset)

			continue
		}

		doesExists := false

		for _, group := range config.Groups {
			if groupName == group.Name {
				doesExists = true

				break
			}
		}

		if doesExists {
			moveCursorUp(2)
			clearLines()
			fmt.Printf("%sError: Group already exists!%s\n%sEnter new group name: %s", colors.Yellow, colors.Reset, colors.Green, colors.Reset)

			continue
		}

		err := fs.CreateGroup(groupName)
		if err != nil {
			moveCursorUp(2)
			clearLines()

			return fmt.Errorf("Create group error: %v", err)
		}

		moveCursorUp(2)
		clearLines()
		fmt.Printf("%sSuccessfully created group!%s\n", colors.Blue, colors.Reset)

		return nil
	}

	return nil
}

func RenderGroupsMenu() error {
	config, err := fs.GetConfig()
	if err != nil {
		return fmt.Errorf("Could not read config file")
	}

	if len(config.Groups) == 0 {
		fmt.Printf("%sNo groups found. Create one first.%s\n", colors.Yellow, colors.Reset)
		return nil
	}

	return nil

}
