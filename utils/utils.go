package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func GetUniqueID() string {
	id := uuid.New().String()
	return strings.Replace(id, "-", "", -1)
}

func ParseBoolean(value string) bool {
	switch strings.ToLower(value) {
	case "false", "0":
		return false
	default:
		return true
	}
}

func ParseInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return value
}
