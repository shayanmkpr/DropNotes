package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type TodoApp struct {
	App   fyne.App
	Todos []string
}

func (t *TodoApp) CreateMenu() *fyne.Menu {
	var items []*fyne.MenuItem

	for i, todo := range t.Todos {
		todoText := todo
		index := i
		item := fyne.NewMenuItem(todoText, func() {
			t.RemoveTodo(index)
		})
		items = append(items, item)
	}

	if len(t.Todos) > 0 {
		items = append(items, fyne.NewMenuItemSeparator())
	}

	addItem := fyne.NewMenuItem("Add Todo...", func() {
		t.ShowAddWindow()
	})
	items = append(items, addItem)

	items = append(items, fyne.NewMenuItemSeparator())
	quitItem := fyne.NewMenuItem("Quit", func() {
		t.App.Quit()
	})
	items = append(items, quitItem)

	return fyne.NewMenu("Todos", items...)
}

func (t *TodoApp) UpdateSystemTray() {
	if desk, ok := t.App.(desktop.App); ok {
		menu := t.CreateMenu()
		desk.SetSystemTrayMenu(menu)
	}
}

func (t *TodoApp) AddTodo(text string) {
	if text != "" {
		t.Todos = append(t.Todos, text)
		t.UpdateSystemTray()
	}
}

func (t *TodoApp) RemoveTodo(index int) {
	if index >= 0 && index < len(t.Todos) {
		t.Todos = append(t.Todos[:index], t.Todos[index+1:]...)
		t.UpdateSystemTray()
	}
}

func (t *TodoApp) ShowAddWindow() {
	w := t.App.NewWindow("Add Todo")
	w.Resize(fyne.NewSize(350, 120))

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter your todo item...")

	addBtn := widget.NewButton("Add", func() {
		t.AddTodo(entry.Text)
		w.Close()
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		w.Close()
	})

	entry.OnSubmitted = func(text string) {
		t.AddTodo(text)
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

