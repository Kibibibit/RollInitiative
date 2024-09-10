package ui

import (
	"fmt"
	"windmills/roll_initiative/models"

	"github.com/awesome-gocui/gocui"
)

const (
	keybindBufferLength = 2
	keybindLeader       = gocui.KeySpace
)

type RootWidget struct {
	view          *gocui.View
	name          string
	x, y          int
	w, h          int
	dataStore     *models.DataStore
	colors        *ColorPalette
	keybindBuffer string
	leaderPressed bool
}

func NewRootWidget(dataStore *models.DataStore, colors *ColorPalette) *RootWidget {
	return &RootWidget{
		name:          NameRootWidget,
		dataStore:     dataStore,
		x:             -1,
		y:             -1,
		keybindBuffer: "",
		leaderPressed: false,
		colors:        colors,
	}
}

func (w *RootWidget) Layout(g *gocui.Gui) error {

	width, height := g.Size()

	w.w = width - 1
	w.h = height - 1

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		w.createKeybinds(g)
		view.Frame = false
		view.BgColor = w.colors.BgColor.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.SetWritePos(0, w.h)

		fmt.Fprintf(view, "Roll Initiative")

	} else if err != nil {
		return err
	}

	return nil

}

func (w *RootWidget) createKeybinds(g *gocui.Gui) error {

	if err := g.SetKeybinding(w.name, keybindLeader, gocui.ModNone, w.onLeader); err != nil {
		return err
	}

	return nil

}

func (w *RootWidget) onLeader(g *gocui.Gui, v *gocui.View) error {
	shortcutView, err := g.View(NameShortcutsWidget)

	if err != nil {
		return err
	}

	shortcutView.Visible = true

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(shortcutView.Name())
		return err
	})

	return nil
}
