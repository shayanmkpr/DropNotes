package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	// "fyne.io/fyne/v2/storage"
)

type NoteApp struct {
	App   fyne.App
	Notes []string
}

const LineSizeThreshold = 100

// Get the file path for storing notes
func getNotesFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".drop_notes.json")
}

// Load notes from file
func (t *NoteApp) LoadNotes() error {
	filePath := getNotesFilePath()
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, start with empty notes
		t.Notes = []string{}
		return nil
	}
	
	// Read file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	// Unmarshal JSON
	return json.Unmarshal(data, &t.Notes)
}

// Save notes to file
func (t *NoteApp) SaveNotes() error {
	filePath := getNotesFilePath()
	
	// Marshal to JSON
	data, err := json.MarshalIndent(t.Notes, "", "  ")
	if err != nil {
		return err
	}
	
	// Write to file
	return ioutil.WriteFile(filePath, data, 0644)
}

func (t *NoteApp) CreateMenu() *fyne.Menu {
	var items []*fyne.MenuItem

	for i, note := range t.Notes {
		noteText := note
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

	// Add quick option to add from clipboard
	clipboardItem := fyne.NewMenuItem("Add from Clipboard", func() {
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		clipText := clipboard.Content()
		if clipText != "" {
			t.AddNote(clipText)
		}
	})
	items = append(items, clipboardItem)

	items = append(items, fyne.NewMenuItemSeparator())
	quitItem := fyne.NewMenuItem("Quit", func() {
		// Use os.Exit instead of App.Quit() for system tray apps
		os.Exit(0)
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
				if i+LineSizeThreshold <= len(text) {
					t.Notes = append(t.Notes, text[i:i+LineSizeThreshold])
				} else {
					t.Notes = append(t.Notes, text[i:])
				}
			}
		} else {
			t.Notes = append(t.Notes, text)
		}
		t.SaveNotes() // Save after adding
		t.UpdateSystemTray()
	}
}

func (t *NoteApp) RemoveNote(index int) {
	// Fixed: should be >= 0, not > 0
	if index >= 0 && index < len(t.Notes) {
		t.Notes = append(t.Notes[:index], t.Notes[index+1:]...)
		t.SaveNotes() // Save after removing
		t.UpdateSystemTray()
	}
}

func (t *NoteApp) ShowAddWindow() {
	// Create a very small, minimal window
	w := t.App.NewWindow("")
	w.Resize(fyne.NewSize(300, 80))
	w.SetFixedSize(true)

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Type note and press Enter...")

	entry.OnSubmitted = func(text string) {
		t.AddNote(text)
		w.Close()
	}

	// Just the entry field, no buttons - cleaner
	content := container.NewVBox(entry)

	w.SetContent(content)
	w.CenterOnScreen()
	w.Show()

	w.Canvas().Focus(entry)
	
	// Auto-close if user clicks away (focus lost)
	w.SetOnClosed(func() {
		// Window closed
	})
}
