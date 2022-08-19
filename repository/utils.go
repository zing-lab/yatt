package repository

import (
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func init() {
	if path, err := os.UserHomeDir(); err == nil {
		filePath = path + "/.yatt/"
	}
}

var (
	once     sync.Once
	lStorage *localStorageRepo
	appName  = "YATT"
	filePath = "~/.yatt/"
	fileName = "storage.xlsx"
	rowLimit = 20
)

const (
	noteSheet   = "note"
	configSheet = "config"
)

var configDetails = map[string]map[string]string{
	"currentRow": {
		"default": "0",
		"row":     "2",
	},
	"currentNoteSheet": {
		"default": "0",
		"row":     "3",
	},
}

func getUniqueID() string {
	id := uuid.New().String()
	return strings.Replace(id, "-", "", -1)
}
