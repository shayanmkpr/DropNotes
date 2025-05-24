package main

import (
	"drop_notes/utils"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
	myApp := app.New()
	myApp.SetIcon(nil)

	todoApp := &utils.TodoApp{
		App:   myApp,
		Todos: []string{"Sample Todo - Click to remove"},
	}

	if desk, ok := myApp.(desktop.App); ok {
		menu := todoApp.CreateMenu()
		desk.SetSystemTrayMenu(menu)
	}

	myApp.Run()
}

