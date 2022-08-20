package repository

import (
	"os"
	"sync"
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
	"marked_only": {
		"default": "0",
		"row":     "4",
	},
}
