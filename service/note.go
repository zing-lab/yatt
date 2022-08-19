package service

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/Kimbbakar/yatt/repository"
)

var (
	// key - id - date - note - description - deleted
	prefixIndent2 = "  "
	prefixIndent4 = "    "
	lineDevider   = "|yatt@yatt|"
)

const (
	KEY     = 0
	ID      = 1
	DATE    = 2
	NOTE    = 3
	DESC    = 4
	DELETED = 5
)

type Note struct {
	id          string
	note        string
	description string
	date        time.Time
	deleted     bool
}

func (n Note) GetID() string {
	return n.id
}

func (n Note) GetDescription() string {
	return n.description
}

func (n Note) String() string {
	resp := "[ - ] "
	if n.deleted {
		resp = "[ √ ]"
	}
	return fmt.Sprintf("%s %s (%s)", resp, n.note, humanize.Time(n.date))
}

type NoteService struct {
}

func (n *NoteService) CreateCommand(note, description string) error {
	repo := repository.GetNewLocalStorage()
	return repo.AddNote(note, description)
}

func (n *NoteService) ListCommand(curPage int) []Note {
	limit := 10
	start, end, limit := curPage*limit, (curPage+1)*limit, limit

	noteList := []Note{}
	repo := repository.GetNewLocalStorage()
	curSheet, err := repo.NextSheet("")
	if err != nil {
		response(err.Error(), true, false, true)
	}
	for c := 0; !(limit <= 0 || curSheet == ""); {
		notes, err := repo.ListNotes(curSheet)
		if err != nil {
			response(err.Error(), true, false, true)
		}

		for i := len(notes) - 1; i >= 0 && limit > 0; i-- {
			deleted, err := strconv.Atoi(notes[i][DELETED])
			if err != nil {
				response(err.Error(), true, false, true)
			}

			if start <= c+i && c+i <= end {
				createdAt, _ := time.Parse(time.RFC1123, notes[i][DATE])
				noteList = append(noteList, Note{
					id:          notes[i][ID],
					note:        notes[i][NOTE],
					deleted:     deleted == 1,
					description: notes[i][DESC],
					date:        createdAt,
				})
				limit--
			}

		}

		c += len(notes)
		curSheet, err = repo.NextSheet(curSheet)
		if err != nil {
			response(err.Error(), true, false, true)
		}
	}

	return noteList
}

func (n *NoteService) FlashStorageCommand() error {
	repo := repository.GetNewLocalStorage()
	repo.FlashStorage()
	return nil
}

func (n *NoteService) ToggleCommand(id string) error {
	repo := repository.GetNewLocalStorage()
	curSheet, err := repo.NextSheet("")
	if err != nil {
		response(err.Error(), true, false, true)
	}

	for {
		if curSheet == "" {
			break
		}

		notes, err := repo.ListNotes(curSheet)
		if err != nil {
			response(err.Error(), true, false, true)
		}

		for i := len(notes) - 1; i >= 0; i-- {
			if strings.HasPrefix(notes[i][ID], id) {
				deleted, _ := strconv.Atoi(notes[i][DELETED])
				row := strings.Split(notes[i][KEY], "-")[2]

				updateValue := make([]interface{}, len(notes[i]))
				for idx, v := range notes[i] {
					updateValue[idx] = v
				}

				updateValue[DELETED] = 1 - deleted
				repo.UpdateNote(curSheet, row, updateValue)
				response("Note has been updated successfully", false, false, true)
				return nil
			}
		}

		curSheet, err = repo.NextSheet(curSheet)
		if err != nil {
			response(err.Error(), true, false, true)
		}
	}

	response("No note found with given ID", false, false, true)
	return nil
}

func (n *NoteService) inputDescription() (string, error) {
	fmt.Println("\nAdd the description[entry empty line to terminate]")
	in := bufio.NewReader(os.Stdin)
	details := ""
	for {
		fmt.Print(prefixIndent4)
		str, err := in.ReadString('\n')
		str = strings.Trim(str, " ")

		if err != nil {
			return "", err
		} else if str == "\n" {
			break
		}

		if len(details) > 0 {
			details += lineDevider
		}
		details += str
	}

	return details, nil
}
