package main

import (
	// "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"drop_notes/utils" // Adjust this path to match your project structure
)

func main() {
	myApp := app.New()
	myApp.SetIcon(nil)
	
	noteApp := &utils.NoteApp{
		App:   myApp,
		Notes: []string{},
	}
	
	// Load existing notes from file
	err := noteApp.LoadNotes()
	if err != nil {
		// If loading fails, just start with empty notes
		noteApp.Notes = []string{}
	}
	
	// Set up system tray
	if desk, ok := myApp.(desktop.App); ok {
		menu := noteApp.CreateMenu()
		desk.SetSystemTrayMenu(menu)
	}
	
	myApp.Run()
}
