package utils

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// NoteApp represents the main application structure
type NoteApp struct {
	App   fyne.App    // The Fyne application instance
	Notes []string    // Slice containing all notes
}

const LineSizeThreshold = 100 // Maximum length of a single note before splitting

// HandleNotes manages all note operations (loading, saving, adding, removing)
func (t *NoteApp) HandleNotes(operation string, noteText string, index int) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
				fmt.Println(err)
		}
		filePath := filepath.Join(homeDir, ".drop_notes.json")

    // filePath := filepath.Join(os.UserHomeDir(), ".drop_notes.json")
    switch operation {
    case "load":
        // If file doesn't exist, start with empty notes
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            t.Notes = []string{}
            return nil
        }
        
        // Read and parse existing notes
        data, err := ioutil.ReadFile(filePath)
        if err != nil {
            return err
        }
        return json.Unmarshal(data, &t.Notes)
        
    case "save":
        // Save notes to file in JSON format
        data, err := json.MarshalIndent(t.Notes, "", "  ")
        if err != nil {
            return err
        }
        return ioutil.WriteFile(filePath, data, 0644)
        
    case "add":
        // Add new note, splitting if it exceeds threshold
        if noteText != "" {
            if len(noteText) > LineSizeThreshold {
                // Split long notes into smaller chunks
                for i := 0; i < len(noteText); i += LineSizeThreshold {
                    end := i + LineSizeThreshold
                    if end > len(noteText) {
                        end = len(noteText)
                    }
                    t.Notes = append(t.Notes, noteText[i:end])
                }
            } else {
                t.Notes = append(t.Notes, noteText)
            }
            t.HandleNotes("save", "", 0)
            t.UpdateUI()
        }
        
    case "remove":
        // Remove note at specified index
        if index >= 0 && index < len(t.Notes) {
            t.Notes = append(t.Notes[:index], t.Notes[index+1:]...)
            t.HandleNotes("save", "", 0)
            t.UpdateUI()
        }
    }
    return nil
}

// UpdateUI handles all UI-related operations (system tray menu and add window)
func (t *NoteApp) UpdateUI() {
    // Create menu items for each note
    var menuItems []*fyne.MenuItem
    
    // Add existing notes
    for i, note := range t.Notes {
        noteIndex := i
        menuItems = append(menuItems, fyne.NewMenuItem(note, func() {
            t.HandleNotes("remove", "", noteIndex)
        }))
    }
    
    // Add separator if there are notes
    if len(t.Notes) > 0 {
        menuItems = append(menuItems, fyne.NewMenuItemSeparator())
    }
    
    // Add utility menu items
    menuItems = append(menuItems,
        fyne.NewMenuItem("Add Note...", t.showAddNoteWindow),
        fyne.NewMenuItem("Add from Clipboard", func() {
            if clipText := t.App.Driver().AllWindows()[0].Clipboard().Content(); clipText != "" {
                t.HandleNotes("add", clipText, 0)
            }
        }),
        fyne.NewMenuItemSeparator(),
        fyne.NewMenuItem("Quit", func() { os.Exit(0) }),
    )
    
    // Update system tray if available
    if desk, ok := t.App.(desktop.App); ok {
        desk.SetSystemTrayMenu(fyne.NewMenu("Notes", menuItems...))
    }
}

// showAddNoteWindow creates and displays a window for adding new notes
func (t *NoteApp) showAddNoteWindow() {
    w := t.App.NewWindow("Add Note")
    w.Resize(fyne.NewSize(300, 80))
    w.SetFixedSize(true)
    
    entry := widget.NewEntry()
    entry.SetPlaceHolder("")
    entry.OnSubmitted = func(text string) {
        t.HandleNotes("add", text, 0)
        w.Close()
    }
    
    w.SetContent(container.NewVBox(entry))
    w.CenterOnScreen()
    w.Show()
    w.Canvas().Focus(entry)
}
