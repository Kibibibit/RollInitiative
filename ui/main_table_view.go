package ui

import (
	"github.com/awesome-gocui/gocui"
)

const MAIN_TABLE_NAME = "main_table"

type MainTableWidget struct {
	name string
	x, y int
}

func NewMainTableWidget(x, y int) *MainTableWidget {

	return &MainTableWidget{
		name: MAIN_TABLE_NAME,
		x:    x,
		y:    y,
	}
}

func (w *MainTableWidget) Layout(g *gocui.Gui) error {

	width, height := g.Size()

	view, err := g.SetView(w.name, w.x, w.y, w.x+width-1, w.y+height-1, 0)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		g.SetCurrentView(w.name)
	}

	view.Title = "Roll Initiative"

	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
