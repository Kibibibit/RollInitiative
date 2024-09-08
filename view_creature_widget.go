package main

import (
	"windmills/roll_initiative/models"

	"github.com/awesome-gocui/gocui"
)

const VIEW_CREATURE_NAME = "view_creature"

type ViewCreatureWidget struct {
	name       string
	w, h       int
	creature   models.Creature
	lastWindow string
}

func NewViewCreatureWidget(name string, w, h int, creature models.Creature, lastWindow string) *ViewCreatureWidget {

	return &ViewCreatureWidget{
		name:       name,
		w:          w,
		h:          h,
		creature:   creature,
		lastWindow: lastWindow,
	}
}

func (w *ViewCreatureWidget) Layout(g *gocui.Gui) error {

	cols, lines := g.Size()

	x := int(cols / 2)
	y := int(lines / 2)

	width := int(w.w / 2)
	height := int(w.h / 2)

	x0 := int(x - width)
	x1 := int(x + width)
	y0 := int(y - height)
	y1 := int(y + height)

	view, err := g.SetView(w.name, x0, y0, x1, y1, 0)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

	}

	view.Title = w.creature.Name

	return nil
}

func switchToViewCreatureWindow(g *gocui.Gui, v *gocui.View, widget *ViewCreatureWidget) error {

	widget.Layout(g)

	_, err := g.SetCurrentView(VIEW_CREATURE_NAME)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	return nil
}
