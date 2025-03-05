package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

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
		SetColumns(0, 30).
		SetBorders(true).
		AddItem(newPrimitive("Debug"), 0, 0, 1, 1, 3, 3, false).
		AddItem(newPrimitive("State"), 0, 1, 1, 1, 3, 3, false)

	// table := tview.NewTable()
	// table.
	// grid.AddItem(table, 0, 0, 1, 1, 3, 3, false)

	app := tview.NewApplication()
	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}

	return &uiState{
		app:  app,
		grid: grid,
	}
}

func (ui *uiState) Debug(args ...any) {
	t := tview.NewTextView()
	fmt.Fprintln(t, args...)
	ui.grid.AddItem(t, 0, 0, 1, 1, 3, 3, false)
}

func (ui *uiState) Error(args ...any) {
	ui.Error(args...)
}
