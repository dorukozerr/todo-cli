package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"github.com/dorukozerr/todo-cli/internal/config"
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

func GenerateNextTodoID(c types.Config) string {
	if len(c.Todos) == 0 {
		return "1"
	}

	maxID := 0
	for _, todo := range c.Todos {
		if id, err := strconv.Atoi(todo.ID); err == nil {
			if id > maxID {
				maxID = id
			}
		}
	}

	return strconv.Itoa(maxID + 1)
}

func CountTodosInGroup(c types.Config, groupName string) int {
	todoCount := 0

	for _, todo := range c.Todos {
		if todo.Group == groupName {
			todoCount = todoCount + 1
		}
	}

	return todoCount
}

func GetUrgencyDisplay(urgency int) (string, string) {
	switch urgency {
	case 5:
		return "CRITICAL", config.Red + config.Bold
	case 4:
		return "HIGH", config.Red
	case 3:
		return "MEDIUM", config.Yellow
	case 2:
		return "LOW", config.Green
	case 1:
		return "MINIMAL", config.Cyan
	default:
		return "UNKNOWN", config.White
	}
}
