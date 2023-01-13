package repository

import (
	"os"
	"sync"

	. "github.com/zing-lab/yatt/utils"
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

var configDetails = map[ConfigKey]map[string]string{
	CurrentRow: {
		"default": "0",
		"row":     "2",
	},
	CurrentNoteSheet: {
		"default": "0",
		"row":     "3",
	},
	MarkedOnly: {
		"default": "0",
		"row":     "4",
	},
	PerPage: {
		"default": "10",
		"row":     "5",
	},
	Tags: {
		"default": "All,ToDo",
		"row":     "6",
	},
	CurrentTagIdx: {
		"default": "0",
		"row":     "7",
	},
}
