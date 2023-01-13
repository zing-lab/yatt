package repository

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zing-lab/yatt/utils"
)

type localStorageRepo struct {
	client *excelize.File
}

func getStorage() *excelize.File {
	f, err := excelize.OpenFile(filePath + fileName)
	if err != nil {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		f = excelize.NewFile()
		if err := f.SaveAs(filePath + fileName); err != nil {
			log.Fatal(err)
		}
	}

	f.NewSheet(configSheet)
	return f
}

func GetNewLocalStorage() *localStorageRepo {
	once.Do(func() {
		lStorage = &localStorageRepo{getStorage()}
	})

	return lStorage
}

func (l *localStorageRepo) GetConfig(key utils.ConfigKey) string {
	v, err := l.client.GetCellValue(configSheet, "B"+configDetails[key]["row"])
	if err != nil || v == "" {
		v = configDetails[key]["default"]
	}

	return v
}

func (l *localStorageRepo) SetConfig(key utils.ConfigKey, value interface{}) error {
	if err := l.client.SetCellValue(configSheet, "A"+configDetails[key]["row"], key); err != nil {
		return err
	}

	if err := l.client.SetCellValue(configSheet, "B"+configDetails[key]["row"], value); err != nil {
		return err
	}

	return l.client.Save()
}

func (l *localStorageRepo) getNewRow() (string, error) {
	curRow, err := strconv.Atoi(l.GetConfig(utils.CurrentRow))
	if err != nil {
		return "", err
	}

	l.SetConfig(utils.CurrentRow, curRow+1)
	return "A" + strconv.Itoa(curRow+1), nil
}

func (l *localStorageRepo) getNoteSheet() (string, error) {
	curRow, err := strconv.Atoi(l.GetConfig(utils.CurrentRow))
	if err != nil {
		return "", err
	}

	curSheet, err := strconv.Atoi(l.GetConfig(utils.CurrentNoteSheet))
	if err != nil {
		return "", err
	}

	if curRow >= rowLimit {
		curSheet++
		l.SetConfig(utils.CurrentRow, 0)
	}

	sheet := noteSheet + "-" + strconv.Itoa(curSheet)
	l.client.NewSheet(sheet)
	l.SetConfig(utils.CurrentNoteSheet, curSheet)

	return sheet, nil
}

func (l *localStorageRepo) AddNote(note, description string) error {
	sheet, err := l.getNoteSheet()
	if err != nil {
		return err
	}
	row, err := l.getNewRow()
	if err != nil {
		return err
	}
	key := appName + "-" + strings.Split(sheet, "-")[1] + "-" + row
	id := utils.GetUniqueID()
	date := time.Now().Format(time.RFC1123)

	// key - id - date - note - description - deleted
	if err := l.client.SetSheetRow(sheet, row, &[]interface{}{key, id, date, note, description, false}); err != nil {
		return err
	}

	return l.client.Save()
}

func (l *localStorageRepo) UpdateNote(sheet, row string, value []interface{}) error {
	if err := l.client.SetSheetRow(sheet, row, &value); err != nil {
		return err
	}

	return l.client.Save()
}

func (l *localStorageRepo) ListNotes(sheetName string) ([][]string, error) {
	return l.client.GetRows(sheetName)
}

func (l *localStorageRepo) FlushStorage() error {
	return os.RemoveAll(filePath)
}

func (l *localStorageRepo) NextSheet(sheetName string) (string, error) {
	if sheetName == "" {
		return l.getNoteSheet()
	}

	data := strings.Split(sheetName, "-")
	if data[1] == "0" {
		return "", nil
	}

	curSheet, err := strconv.Atoi(data[1])
	if err != nil {
		return "", err
	}
	return data[0] + "-" + strconv.Itoa(curSheet-1), nil
}

func (l *localStorageRepo) GetTags() []string {
	tags := l.GetConfig(utils.Tags)
	return strings.Split(tags, ",")
}

func (l *localStorageRepo) GetCurrentTagIndex() int {
	idxStr := l.GetConfig(utils.CurrentTagIdx)
	idx, _ := strconv.Atoi(idxStr)
	return idx
}
