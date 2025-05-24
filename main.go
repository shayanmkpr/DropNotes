package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/driver/desktop"
)

type TodoApp struct {
	app   fyne.App
	todos []string
}

func main() {
	myApp := app.New()
	myApp.SetIcon(nil)
	
	todoApp := &TodoApp{
		app:   myApp,
		todos: []string{"Sample Todo - Click to remove"},
	}
	
	if desk, ok := myApp.(desktop.App); ok {
		menu := todoApp.createMenu()
		desk.SetSystemTrayMenu(menu)
	}
	
	myApp.Run()
}

func (t *TodoApp) createMenu() *fyne.Menu {
	var items []*fyne.MenuItem
	
	for i, todo := range t.todos {
		todoText := todo
		index := i
		item := fyne.NewMenuItem(todoText, func() {
			t.removeTodo(index)
		})
		items = append(items, item)
	}
	
	if len(t.todos) > 0 {
		items = append(items, fyne.NewMenuItemSeparator())
	}
	
	addItem := fyne.NewMenuItem("Add Todo...", func() {
		t.showAddWindow()
	})
	items = append(items, addItem)
	
	items = append(items, fyne.NewMenuItemSeparator())
	quitItem := fyne.NewMenuItem("Quit", func() {
		t.app.Quit()
	})
	items = append(items, quitItem)
	
	return fyne.NewMenu("Todos", items...)
}

func (t *TodoApp) updateSystemTray() {
	if desk, ok := t.app.(desktop.App); ok {
		menu := t.createMenu()
		desk.SetSystemTrayMenu(menu)
	}
}

func (t *TodoApp) addTodo(text string) {
	if text != "" {
		t.todos = append(t.todos, text)
		t.updateSystemTray()
	}
}

func (t *TodoApp) removeTodo(index int) {
	if index >= 0 && index < len(t.todos) {
		t.todos = append(t.todos[:index], t.todos[index+1:]...)
		t.updateSystemTray()
	}
}

func (t *TodoApp) showAddWindow() {
	w := t.app.NewWindow("Add Todo")
	w.Resize(fyne.NewSize(350, 120))
	
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter your todo item...")
	
	addBtn := widget.NewButton("Add", func() {
		t.addTodo(entry.Text)
		w.Close()
	})
	
	cancelBtn := widget.NewButton("Cancel", func() {
		w.Close()
	})
	
	entry.OnSubmitted = func(text string) {
		t.addTodo(text)
		w.Close()
	}
	
	buttons := container.NewHBox(addBtn, cancelBtn)
	content := container.NewVBox(
		widget.NewLabel("Add a new todo:"),
		entry,
		buttons,
	)
	
	w.SetContent(content)
	w.CenterOnScreen()
	w.Show()
	
	w.Canvas().Focus(entry)
}
