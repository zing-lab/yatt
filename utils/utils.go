package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GetUniqueID() string {
	id := uuid.New().String()
	return strings.Replace(id, "-", "", -1)
}

func ParseBoolean(value string) int {
	switch strings.ToLower(value) {
	case "false", "0":
		return 0
	default:
		return 1
	}
}
