package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dorukozerr/todo-cli/internal/types"
	"github.com/dorukozerr/todo-cli/internal/utils"
)

func GetConfig() (*types.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".config", "todo-cli")
	configPath := filepath.Join(configDir, "config.json")

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		emptyConfig := &types.Config{
			Groups:      []types.Group{},
			ActiveGroup: "",
			Todos:       []types.Todo{},
		}

		configData, err := json.MarshalIndent(emptyConfig, "", "  ")
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(configPath, configData, 0644)
		if err != nil {
			return nil, err
		}

		return emptyConfig, nil
	}

	if err != nil {
		return nil, err
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config types.Config

	err = json.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *types.Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "todo-cli")
	configPath := filepath.Join(configDir, "config.json")

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, configData, 0644)
}

func CreateGroup(name string) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	for _, group := range config.Groups {
		if group.Name == name {
			return fmt.Errorf("group '%s' already exists", name)
		}
	}

	id, err := utils.GenerateRandomID()
	if err != nil {
		return err
	}

	newGroup := types.Group{
		ID:   id,
		Name: name,
	}

	config.Groups = append(config.Groups, newGroup)
	config.ActiveGroup = newGroup.Name

	return SaveConfig(config)
}

func DeleteGroup(groupID string) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	var newGroups []types.Group
	var deletedGroupName string

	for _, group := range config.Groups {
		if group.ID != groupID {
			newGroups = append(newGroups, group)
		} else {
			deletedGroupName = group.Name
		}
	}

	if deletedGroupName == "" {
		return fmt.Errorf("group not found")
	}

	config.Groups = newGroups

	var newTodos []types.Todo

	for _, todo := range config.Todos {
		if todo.Group != deletedGroupName {
			newTodos = append(newTodos, todo)
		}
	}

	config.Todos = newTodos

	if config.ActiveGroup == deletedGroupName {
		config.ActiveGroup = ""
	}

	return SaveConfig(config)
}

func SetActiveGroup(groupName string) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	found := false

	for _, group := range config.Groups {
		if group.Name == groupName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("group '%s' does not exist", groupName)
	}

	config.ActiveGroup = groupName

	return SaveConfig(config)
}
