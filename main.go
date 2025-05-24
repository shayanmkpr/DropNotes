package main

import (
	"drop_notes/utils"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
  "fyne.io/fyne/v2/theme"
)

func main() {
	myApp := app.New()
	myApp.SetIcon(theme.Icon(theme.IconNameConfirm))

	todoApp := &utils.NoteApp{
		App:   myApp,
		Notes: []string{"Drop Down Notes"},
	}

	if desk, ok := myApp.(desktop.App); ok {
		menu := todoApp.CreateMenu()
		desk.SetSystemTrayMenu(menu)
	}

	myApp.Run()
}

