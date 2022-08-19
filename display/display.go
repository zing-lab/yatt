package display

import (
	"strings"

	"github.com/Kimbbakar/yatt/service"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	curPage int
	app     *tview.Application
	list    *tview.List
	srvc    service.NoteService
	noteId  map[int]string
)

func init() {
	initDisplay()
}

func initDisplay() {
	srvc = service.NoteService{}
	app = tview.NewApplication()
	list = tview.NewList()
	list.SetBorder(true).SetTitle("YATT [Help = ctrl + h]").SetTitleAlign(tview.AlignLeft)
	renderListCommand()
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			idx := list.GetCurrentItem()
			id := noteId[idx]
			srvc.ToggleCommand(id)
			renderListCommand()
		case tcell.KeyCtrlI:
			createNoteCommand()
		case tcell.KeyCtrlO:
			if curPage > 0 {
				curPage--
			}
			renderListCommand()
		case tcell.KeyCtrlP:
			curPage++
			renderListCommand()
		case tcell.KeyCtrlF:
			flashCommnad()
		case tcell.KeyCtrlH:
			helpCommnad()
		}
		return event
	})

	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

func createNoteCommand() {
	form := tview.NewForm().
		AddInputField("Note", "", 30, nil, nil).
		AddInputField("Description", "", 50, nil, nil)
	app.SetRoot(form, true).SetFocus(form)

	form = form.AddButton("Save", func() {
		note := form.GetFormItemByLabel("Note").(*tview.InputField).GetText()
		if note = strings.TrimSpace(note); note == "" {
			return
		}

		description := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
		srvc.CreateCommand(note, description)
		renderListCommand()
		app.SetRoot(list, true).SetFocus(list)
	}).AddButton("Quit", func() {
		app.SetRoot(list, true).SetFocus(list)
	})
}

func flashCommnad() {
	modal := tview.NewModal().
		SetText("Do you want to flash all notes?").
		AddButtons([]string{"Yes", "No"}).
		SetButtonBackgroundColor(tcell.ColorDarkRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				srvc.FlashStorageCommand()
				app.Stop()
			}
			app.SetRoot(list, true).SetFocus(list)
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func helpCommnad() {
	modal := tview.NewModal().
		SetText("Shortcut \n mark/unmark = enter \n new note = ctrl + i \n previous page = ctrl + o \n next page = ctrl + p \n flash = ctrl + f").
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(list, true).SetFocus(list)
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func renderListCommand() {
	noteId = map[int]string{}
	idx := list.GetCurrentItem()

	list.Clear()
	notes := srvc.ListCommand(curPage)
	for idx, note := range notes {
		noteId[idx] = note.GetID()
		list = list.AddItem(note.String(), note.GetDescription(), rune(0), nil)
	}

	list.SetCurrentItem(idx)
}
