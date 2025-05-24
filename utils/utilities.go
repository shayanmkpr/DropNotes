package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type NoteApp struct {
	App   fyne.App
	Notes []string
}

const LineSizeThreshold = 100

func (t *NoteApp) CreateMenu() *fyne.Menu {
	var items []*fyne.MenuItem

	for i, note := range t.Notes {
		noteText :=note 
		index := i
		item := fyne.NewMenuItem(noteText, func() {
			t.RemoveNote(index)
		})
		items = append(items, item)
	}

	if len(t.Notes) > 0 {
		items = append(items, fyne.NewMenuItemSeparator())
	}

	addItem := fyne.NewMenuItem("Add Note...", func() {
		t.ShowAddWindow()
	})
	items = append(items, addItem)

	items = append(items, fyne.NewMenuItemSeparator())
	quitItem := fyne.NewMenuItem("Quit", func() {
		t.App.Quit() // There is a problem here.
	})
	items = append(items, quitItem)

	return fyne.NewMenu("Notes", items...)
}

func (t *NoteApp) UpdateSystemTray() {
	if desk, ok := t.App.(desktop.App); ok {
		menu := t.CreateMenu()
		desk.SetSystemTrayMenu(menu)
	}
}

func (t *NoteApp) AddNote(text string) {
	if text != "" {
    if len(text) > LineSizeThreshold {
      for i := 0; i < len(text); i += LineSizeThreshold {
        if i + LineSizeThreshold <= len(text) {
          t.Notes = append(t.Notes, text[i: i + LineSizeThreshold])
        } else {
          t.Notes = append(t.Notes, text[i:])
        }
      }
    } else {
      t.Notes = append(t.Notes, text)
    }
    t.UpdateSystemTray()
  }
}

func (t *NoteApp) RemoveNote(index int) {
	if index > 0 && index < len(t.Notes) {
		t.Notes = append(t.Notes[:index], t.Notes[index+1:]...)
		t.UpdateSystemTray()
	}
}

func (t *NoteApp) TrimNote(note string) []string {
  var chopped_note []string
  for i := 0; i <= len(note); i += LineSizeThreshold{
    end := i + LineSizeThreshold
    start := i
    if end > len(note){
      end = len(note)
    }
    chopped_note = append(chopped_note, note[start: end])
  }
  return chopped_note
}

func (t *NoteApp) ShowAddWindow() {
	w := t.App.NewWindow("Add")
	w.Resize(fyne.NewSize(350, 120))

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Add Item")

	addBtn := widget.NewButton("Add", func() {
		t.AddNote(entry.Text)
		w.Close()
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		w.Close()
	})

	entry.OnSubmitted = func(text string) {
		t.AddNote(text)
		w.Close()
	}

	buttons := container.NewHBox(addBtn, cancelBtn)
	content := container.NewVBox(
		widget.NewLabel("Add a new note:"),
		entry,
		buttons,
	)

	w.SetContent(content)
	w.CenterOnScreen()
	w.Show()

	w.Canvas().Focus(entry)
}

