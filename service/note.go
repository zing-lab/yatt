package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/zing-lab/yatt/repository"
	"github.com/zing-lab/yatt/utils"
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
	TagIdx  = 6
	MaxLen  = 7
)

const (
	DefaultTagIdx = 0
)

type Note struct {
	id          string
	note        string
	description string
	date        time.Time
	deleted     bool
	tagIndex    int
}

func (n Note) GetID() string {
	return n.id
}

func (n Note) GetDescription() string {
	return n.description
}

func (n Note) GetTag() int {
	return n.tagIndex
}

func (n Note) String() string {
	resp := "[ - ] "
	if n.deleted {
		resp = "[ √ ]"
	}
	return fmt.Sprintf("%s %s (posted %s)", resp, n.note, humanize.Time(n.date))
}

type NoteService struct {
}

func (n *NoteService) SanitizeText(s string) string {
	s = strings.TrimPrefix(s, "[ - ]")
	s = strings.TrimPrefix(s, "[ √ ]")
	if idx := strings.Index(s, "(posted"); idx != -1 {
		s = s[:idx]
	}
	s = strings.TrimSpace(s)
	return s
}

func (n *NoteService) CreateCommand(note, description string, tagID int) error {
	repo := repository.GetNewLocalStorage()
	return repo.AddNote(note, description, tagID)
}

func (n *NoteService) ListCommand(curPage int) []Note {
	repo := repository.GetNewLocalStorage()

	curTagIdx := repo.GetCurrentTagIndex()
	limit := utils.ParseInt(repo.GetConfig("per_page"))
	start, end, limit := curPage*limit, (curPage+1)*limit-1, limit
	markedOnly := utils.ParseBoolean(repo.GetConfig("marked_only"))
	noteList := []Note{}

	curSheet, err := repo.NextSheet("")
	if err != nil {
		response(err.Error(), true, false, true)
	}
	for count := 0; !(limit <= 0 || curSheet == ""); {
		notes, err := repo.ListNotes(curSheet)
		if err != nil {
			response(err.Error(), true, false, true)
		}

		for i := len(notes) - 1; i >= 0 && limit > 0; i, count = i-1, count+1 {
			deleted := utils.ParseBoolean(notes[i][DELETED])
			if err != nil {
				response(err.Error(), true, false, true)
			}

			noteTag := DefaultTagIdx
			if len(notes[i]) > TagIdx {
				noteTag = utils.ParseInt(notes[i][TagIdx])
			}

			if noteTag != curTagIdx || (markedOnly && deleted) {
				start, end = start+1, end+1
				continue
			}

			if !(start <= count && count <= end) {
				continue
			}

			limit--
			createdAt, _ := time.Parse(time.RFC1123, notes[i][DATE])
			noteList = append(noteList, Note{
				id:          notes[i][ID],
				note:        notes[i][NOTE],
				deleted:     deleted,
				description: notes[i][DESC],
				date:        createdAt,
				tagIndex:    noteTag,
			})
		}

		curSheet, err = repo.NextSheet(curSheet)
		if err != nil {
			response(err.Error(), true, false, true)
		}
	}

	return noteList
}

func (n *NoteService) FlushStorageCommand() error {
	repo := repository.GetNewLocalStorage()
	repo.FlushStorage()
	return nil
}

func (n *NoteService) ToggleCommand(id string) error {
	return n.UpdateCommand(id, func(note []interface{}) []interface{} {
		deleted := utils.ParseBoolean(note[DELETED].(string))
		note[DELETED] = !deleted
		return note
	})
}

func (n *NoteService) EditCommand(id, title, description string, tagIdx int) error {
	return n.UpdateCommand(id, func(note []interface{}) []interface{} {
		note[NOTE] = title
		note[DESC] = description
		note[TagIdx] = tagIdx
		return note
	})
}

func (n *NoteService) GetNote(id string) *Note {
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
				noteTag := DefaultTagIdx
				if len(notes[i]) > TagIdx {
					noteTag = utils.ParseInt(notes[i][TagIdx])
				}

				deleted := utils.ParseBoolean(notes[i][DELETED])
				if err != nil {
					response(err.Error(), true, false, true)
				}

				createdAt, _ := time.Parse(time.RFC1123, notes[i][DATE])
				return &Note{
					id:          notes[i][ID],
					note:        notes[i][NOTE],
					deleted:     deleted,
					description: notes[i][DESC],
					date:        createdAt,
					tagIndex:    noteTag,
				}
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

func (n *NoteService) UpdateCommand(id string, updateFunc func([]interface{}) []interface{}) error {
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

				row := strings.Split(notes[i][KEY], "-")[2]
				updateValue := make([]interface{}, MaxLen)
				for idx, v := range notes[i] {
					updateValue[idx] = v
				}

				repo.UpdateNote(curSheet, row, updateFunc(updateValue))
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

func (n *NoteService) GetConfig(key utils.ConfigKey) string {
	repo := repository.GetNewLocalStorage()
	return repo.GetConfig(key)
}

func (n *NoteService) SetConfig(key utils.ConfigKey, value interface{}) error {
	repo := repository.GetNewLocalStorage()
	if err := repo.SetConfig(key, value); err != nil {

		return response("Failed to update setting", true, false, false)
	}

	return nil
}

func (n *NoteService) GetTagDetails() ([]string, int) {
	repo := repository.GetNewLocalStorage()
	return repo.GetTags(), repo.GetCurrentTagIndex()
}

func (n *NoteService) GetTagName() string {
	repo := repository.GetNewLocalStorage()
	tags, idx := repo.GetTags(), repo.GetCurrentTagIndex()
	return strings.Title(tags[idx])
}

func (n *NoteService) AddNewTag(newTag string) int {
	repo := repository.GetNewLocalStorage()
	tags := repo.GetTags()

	tags = append(tags, newTag)
	repo.SetConfig(utils.Tags, strings.Join(tags, ","))
	return len(tags) - 1
}

func (n *NoteService) IsTagValid(newTag string) error {
	repo := repository.GetNewLocalStorage()
	tags := repo.GetTags()

	for _, tag := range tags {
		if strings.EqualFold(tag, newTag) {
			return fmt.Errorf("Tag#(%s) already exist", newTag)
		}
	}

	return nil
}
