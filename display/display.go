package display

import (
	"fmt"
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
			showNoteCommand(nil, "", "")
		case tcell.KeyCtrlE:
			idx := list.GetCurrentItem()
			id := noteId[idx]

			note, description := list.GetItemText(idx)
			showNoteCommand(&id, note, description)
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
			flushCommand()
		case tcell.KeyCtrlS:
			settingCommand()
		case tcell.KeyCtrlH:
			helpCommand()
		case tcell.KeyESC:
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

func showNoteCommand(id *string, note, description string) {
	var noteObj *service.Note
	if id != nil {
		noteObj = srvc.GetNote(*id)
	}

	note = srvc.SanitizeText(note)
	tags, curTagIdx := srvc.GetTagDetails()
	form := tview.NewForm().
		SetFieldBackgroundColor(tcell.ColorAntiqueWhite).
		SetFieldTextColor(tcell.ColorBlack).
		AddInputField("Note", note, 30, nil, nil).
		AddInputField("Description", description, 60, nil, nil)

	if noteObj != nil {
		curTagIdx = noteObj.GetTag()
	}

	form = form.AddDropDown(CurrentTag, tags, curTagIdx, nil)
	app.SetRoot(form, true).SetFocus(form)

	form = form.AddButton(Save, func() {
		note := form.GetFormItemByLabel("Note").(*tview.InputField).GetText()
		if note = strings.TrimSpace(note); note == "" {
			return
		}

		description := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
		curTagIdx, _ = form.GetFormItemByLabel(CurrentTag).(*tview.DropDown).GetCurrentOption()

		if id == nil {
			srvc.CreateCommand(note, description, curTagIdx)
		} else {
			srvc.EditCommand(*id, note, description, curTagIdx)
		}

		renderListCommand()
	})

	form = form.AddButton("Quit", func() {
		app.SetRoot(list, true).SetFocus(list)
	})
}

func flushCommand() {
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

func settingCommand() {
	tags, curTagIdx := srvc.GetTagDetails()
	form := tview.NewForm().
		AddCheckbox(ShowMarkedOnly, utils.ParseBoolean(srvc.GetConfig(utils.MarkedOnly)), nil).
		AddInputField(PerPage, srvc.GetConfig(utils.PerPage), 5, nil, nil).
		AddInputField(NewTag, utils.EmptyString, 20, nil, nil).
		AddDropDown(CurrentTag, tags, curTagIdx, nil)

	form.GetFormItemByLabel(ShowMarkedOnly).(*tview.Checkbox).SetCheckedString("√")
	app.SetRoot(form, true).SetFocus(form)

	form = form.AddButton(Save, func() {
		perPage := form.GetFormItemByLabel(PerPage).(*tview.InputField).GetText()
		value := utils.ParseInt(perPage)
		if value == 0 {
			showModal("Per Page value should be a positive number!", form)
			return
		}

		checkbox := (form.GetFormItemByLabel(ShowMarkedOnly).(*tview.Checkbox))

		srvc.SetConfig(utils.MarkedOnly, checkbox.IsChecked())
		srvc.SetConfig(utils.PerPage, value)
		curTagIdx, _ = form.GetFormItemByLabel(CurrentTag).(*tview.DropDown).GetCurrentOption()
		if newTag := form.GetFormItemByLabel(NewTag).(*tview.InputField).GetText(); newTag != utils.EmptyString {
			if err := srvc.IsTagValid(newTag); err != nil {
				showModal(err.Error(), form)
				return
			}

			curTagIdx = srvc.AddNewTag(newTag)
		}

		srvc.SetConfig(utils.CurrentTagIdx, curTagIdx)

		curPage = 0
		renderListCommand()
	}).AddButton(Cancel, func() {
		app.SetRoot(list, true).SetFocus(list)
	})
}

func helpCommand() {
	s := "Shortcut\n********"
	addShortCutFunc := func(shortC string) {
		s = s + "\n" + shortC
	}

	addShortCutFunc("Mark/Unmark = ENTER")
	addShortCutFunc("New Note = CTRL + I")
	addShortCutFunc("Edit Note = CTRL + E")
	addShortCutFunc("Previous Page = CTRL + O")
	addShortCutFunc("Next Page = CTRL + P")
	addShortCutFunc("Flush = CTRL + F")
	addShortCutFunc("Setting = CTRL + S")

	modal := tview.NewModal().
		SetText(s).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			renderListCommand()
		})
	app.SetRoot(modal, true).SetFocus(modal)
}

func showModal(txt string, backToPage tview.Primitive) {
	modal := tview.NewModal().
		SetText(txt).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(backToPage, true).SetFocus(backToPage)
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

	setTitle()
	list.SetCurrentItem(idx)
	app.SetRoot(list, true).SetFocus(list)
}

func setTitle() {
	title := fmt.Sprintf("YATT [HELP = Ctrl + H] [TAG = %s]", srvc.GetTagName())
	list.SetTitle(title)
}
