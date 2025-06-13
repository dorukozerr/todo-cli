package fs

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/dorukozerr/todo-cli/internal/types"
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
