package ui

import (
	"fmt"
	"log"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

type ViewSpellWidget struct {
	view           *gocui.View
	name           string
	x, y           int
	w, h           int
	colW           int
	colors         *ColorPalette
	dataStore      *models.DataStore
	previousWidget string
	spell          *models.Spell
}

func NewViewSpellWidget(dataStore *models.DataStore, previousWidget string, colors *ColorPalette, spell *models.Spell) *ViewSpellWidget {
	return &ViewSpellWidget{
		name:           NameViewCreatureWidget,
		dataStore:      dataStore,
		colors:         colors,
		previousWidget: previousWidget,
		spell:          spell,
	}

}

func (w *ViewSpellWidget) Layout(g *gocui.Gui) error {

	width, height := g.Size()

	w.w = utils.Clamp(width/2, 70, width-4)
	w.h = height - 5

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

	w.colW = w.w / 3

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		w.createKeybinds(g)
		view.Frame = true
		view.BgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.SelBgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.Title = "View Spell"
		view.Visible = true

	} else if err != nil {
		return err
	}

	view.SetWritePos(1, 1)

	fmt.Fprint(view,
		ApplyBold(
			fmt.Sprintf("%s - Level %d %s ", w.spell.Name, w.spell.Level, w.spell.School),
			w.colors.FgColor,
		),
	)

	return nil
}

func (w *ViewSpellWidget) createKeybinds(g *gocui.Gui) error {

	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.onClose); err != nil {
		return err
	}
	return nil

}

func (w *ViewSpellWidget) onClose(g *gocui.Gui, v *gocui.View) error {

	nextView, err := g.View(NameRootWidget)
	if err != nil {
		log.Println(err)
		return err
	}

	w.view.Visible = false
	g.DeleteKeybindings(w.name)
	g.DeleteView(w.name)
	w.view.Frame = false
	w.view.Clear()

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(nextView.Name())
		return err
	})

	return nil
}
