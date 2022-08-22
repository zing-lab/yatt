package display

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/zing-lab/yatt/service"
	"github.com/zing-lab/yatt/utils"
)

var (
	curPage int
	app     *tview.Application
	list    *tview.List
	srvc    service.NoteService
	noteId  map[int]string
)

func Show() {
	srvc = service.NoteService{}
	app = tview.NewApplication()
	list = tview.NewList()
	list.SetBorder(true).
		SetTitle("YATT [HELP = Ctrl + h]").
		SetTitleAlign(tview.AlignLeft)
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
			if list.GetItemCount() > 0 {
				curPage++
			}
			renderListCommand()
		case tcell.KeyCtrlF:
			flushCommnad()
		case tcell.KeyCtrlS:
			settingCommnad()
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
		SetFieldBackgroundColor(tcell.ColorAntiqueWhite).
		SetFieldTextColor(tcell.ColorBlack).
		AddInputField("Note", "", 30, nil, nil).
		AddInputField("Description", "", 60, nil, nil)
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

func flushCommnad() {
	modal := tview.NewModal().
		SetText("Do you want to flush all notes?").
		AddButtons([]string{"No", "Yes"}).
		SetButtonBackgroundColor(tcell.ColorDarkRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				srvc.FlushStorageCommand()
				app.Stop()
			}
			app.SetRoot(list, true).SetFocus(list)
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func settingCommnad() {
	form := tview.NewForm().
		AddCheckbox("Marked Only", utils.ParseBoolean(srvc.GetConfig("marked_only")), nil).
		AddInputField("Per Page", srvc.GetConfig("per_page"), 5, nil, nil)
	form.GetFormItemByLabel("Marked Only").(*tview.Checkbox).SetCheckedString("√")

	app.SetRoot(form, true).SetFocus(form)

	form = form.AddButton("Save", func() {
		perPage := form.GetFormItemByLabel("Per Page").(*tview.InputField).GetText()
		value := utils.ParseInt(perPage)
		if value == 0 {
			return
		}

		checkbox := (form.GetFormItemByLabel("Marked Only").(*tview.Checkbox))

		srvc.SetConfig("marked_only", checkbox.IsChecked())
		srvc.SetConfig("per_page", value)

		curPage = 0
		renderListCommand()
		app.SetRoot(list, true).SetFocus(list)
	}).AddButton("Cancel", func() {
		app.SetRoot(list, true).SetFocus(list)
	})
}

func helpCommnad() {
	modal := tview.NewModal().
		SetText("Shortcut \n Mark/Unmark = Enter \n New note = Ctrl + i \n Previous page = Ctrl + o \n Next page = Ctrl + p \n Flush = Ctrl + f \n Setting = Ctrl + s").
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

	if idx > len(notes) {
		idx = len(notes) - 1
	}

	list.SetCurrentItem(idx)
}
