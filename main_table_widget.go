package main

import (
	"fmt"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
)

const MAIN_TABLE_NAME = "main_table"

type MainTableWidget struct {
	name string
	x, y int
	data *IniativeTracker
}

func NewMainTableWidget(name string, x, y int) *MainTableWidget {

	return &MainTableWidget{
		name: name,
		x:    x,
		y:    y,
		data: &IniativeTracker{},
	}
}

func AddMainTableWidgetKeybinds(g *gocui.Gui, addCreature *AddCreatureWidget) {
	if err := g.SetKeybinding(MAIN_TABLE_NAME, 'a', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newAddCreatureWindow(g, v, addCreature)
		}); err != nil {
		log.Panicln(err)
		os.Exit(1)
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

		for index, combatant := range w.data.combatants {
			view.SetWritePos(1, index)
			str := fmt.Sprintf("\x1b[37;2m%d\x1b[0m %s (%s)", index+1, combatant.creature.Name, combatant.tag)
			fmt.Fprintf(view, str)
		}

	}

	view.Title = "Roll Initiative"

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
