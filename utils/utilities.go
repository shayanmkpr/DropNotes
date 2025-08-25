package utils

import (
	"fmt"
	// "time"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type MyNote struct {
	Text string
	Status bool
}

type NoteApp struct {
	App   fyne.App    // The Fyne application instance
	Notes []MyNote    // Slice containing all notes
}

const LineSizeThreshold = 100 // Maximum length of a single note before splitting
// HandleNotes manages all note operations (loading, saving, adding, removing)
func (t *NoteApp) HandleNotes(operation string, noteText string, index int) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
			fmt.Println(err)
	}
	folderPath := filepath.Join(homeDir, "/.drop_notes")
	// check if the folder is there. if it wasnt, then make it.
	// everything is going to be there.
	if _, err := os.Stat(folderPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("~/.drop_notes folder does note exist. Creating one.")
			err := os.Mkdir(folderPath, 0755)
			if err != nil {
				println(err)
			} else {
				println("Folder created:", folderPath)
			}
		}
	}
	filePath := filepath.Join(folderPath, ".drop_notes.json")
	// check if the json file is there.
	// if not, then make it so we can write it.
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("~/.drop_notes.json file does note exist. Creating one.")
			_, err := os.Create(filePath)
			if err != nil {
				println(err)
			} else {
				println("File created:", filePath)
			}
		}
	}
    switch operation {
    case "load":

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
					newNote := MyNote{
						Text: noteText[i:end],
						Status: false,
					}
                    t.Notes = append(t.Notes, newNote)
                }
            } else {
				newNote := MyNote{
					Text: noteText,
					Status: false,
				}
                t.Notes = append(t.Notes, newNote)
            }
            t.HandleNotes("save", "", 0)
            t.UpdateUI()
        }
	case "done":
		if index >= 0 && index < len(t.Notes) {
			t.Notes[index].Status = true
			t.HandleNotes("save", "", 0)
			t.UpdateUI()
		}
        
	case "undone":
		if index >= 0 && index < len(t.Notes) {
			t.Notes[index].Status = false
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
    // var doneItem []*fyne.MenuItem
    
    // Add existing notes
    for i, note := range t.Notes {
        noteIndex := i
		// adding the items that are already done:
		if note.Status == true{

			menuItems = append(menuItems, fyne.NewMenuItem("[x]"+note.Text, func() {
				t.HandleNotes("undone", "", noteIndex)
			}))
			// doneItem = append(doneItem, fyne.NewMenuItem("[x]"+note.Text, func(){ // a method to handle the changes that happen after making the check box toggle
			// 	t.HandleNotes("remove", "", noteIndex)
			// 	// t.Notes[noteIndex].Status = !t.Notes[noteIndex].Status // Toggle
			// 	// t.HandleNotes("save", "", 0)
			// 	// t.UpdateUI()
			// }))

		} else {

			menuItems = append(menuItems, fyne.NewMenuItem("[ ]"+note.Text, func() {
				t.HandleNotes("done", "", noteIndex)
			}))

		}
		//
		// menuItems = append(menuItems, doneItem...)

    }
    
    // Add separator if there are notes
    if len(t.Notes) > 0 {
        menuItems = append(menuItems, fyne.NewMenuItemSeparator())
    }
    
    // Add utility menu items
    menuItems = append(menuItems,
        fyne.NewMenuItem("Remove all Done", t.removeDoneItems),
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

func showConfirmWindow(t *NoteApp, callback func(bool)) {
    w := t.App.NewWindow("Confirm")
    w.Resize(fyne.NewSize(200, 100))
    w.SetFixedSize(true)
    label := widget.NewLabel("Are you sure?")
    
    yesButton := widget.NewButton("Yes", func() {
        fmt.Printf("Yes clicked\n")
        callback(true)  // <- This replaces "result = true; return result"
        w.Close()
    })
    
    noButton := widget.NewButton("No", func() {
        fmt.Printf("No clicked\n")
        callback(false) // <- This replaces "result = false; return result"
        w.Close()
    })
    
    w.SetOnClosed(func() {
		callback(false)
    })
    
    w.SetContent(container.NewVBox(label, container.NewHBox(yesButton, noButton)))
    w.CenterOnScreen()
    w.Show()
    
}

func (t *NoteApp) removeDoneItems() {
    // Pass a function that contains "what to do when user responds"
    showConfirmWindow(t, func(confirm bool) {
        // This code runs LATER when user clicks a button
        fmt.Printf("%t 1 \n", confirm)
        
        if confirm == true {
            fmt.Printf("User confirmed, removing items\n")
            for i, note := range(t.Notes){
                if note.Status == true{
                    t.HandleNotes("remove", "", i)
                }
            }
        } else {
            fmt.Printf("%t 2 - User cancelled\n", confirm)
        }
    })
    fmt.Println("Dialog shown, removeDoneItems() function ending")
}
