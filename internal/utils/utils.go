package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/dorukozerr/todo-cli/internal/types"
)

func GenerateRandomID() (string, error) {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func CountTodosInGroup(config types.Config, groupName string) int {
	todoCount := 0

	for _, todo := range config.Todos {
		if todo.Group == groupName {
			todoCount = todoCount + 1
		}
	}

	return todoCount
}
