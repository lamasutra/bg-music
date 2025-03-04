package ui

import "github.com/rivo/tview"

func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}

type uiState struct {
	app  *tview.Application
	grid *tview.Grid
}

func NewTui() *uiState {
	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Map"), 0, 1, 1, 1, 3, 3, false).
		AddItem(newPrimitive("Debug"), 0, 2, 1, 1, 3, 3, false)

	app := tview.NewApplication()
	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}

	return &uiState{
		app:  app,
		grid: grid,
	}
}
