package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"

	"drop_notes/utils"
)

func main() {
	// Create and initialize the application
	myApp := app.New()
	myApp.SetIcon(theme.Icon(theme.IconNameConfirm))

	// Create note app instance
	noteApp := &utils.NoteApp{
		App:   myApp,
		Notes: []string{},
	}

	// Load existing notes or start fresh
	if err := noteApp.HandleNotes("load", "", 0); err != nil {
		noteApp.Notes = []string{}
	}

	// Initialize system tray if available
	if _, ok := myApp.(desktop.App); ok {
		noteApp.UpdateUI()
	}

	// Start the application
	myApp.Run()
}
