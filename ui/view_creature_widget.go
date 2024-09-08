package ui

import (
	"log"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

type ViewCreatureWidget struct {
	view           *gocui.View
	name           string
	x, y           int
	w, h           int
	colors         *ColorPalette
	dataStore      *models.DataStore
	previousWidget string
	creature       *models.Creature
}

func NewViewCreatureWidget(dataStore *models.DataStore, previousWidget string, colors *ColorPalette, creature *models.Creature) *ViewCreatureWidget {
	return &ViewCreatureWidget{
		name:           NameViewCreatureWidget,
		dataStore:      dataStore,
		colors:         colors,
		previousWidget: previousWidget,
		creature:       creature,
	}

}

func (w *ViewCreatureWidget) Layout(g *gocui.Gui) error {

	width, height := g.Size()

	w.w = utils.Clamp(width-50, height-5, width)
	w.h = height - 10

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		w.createKeybinds(g)
		view.Frame = true
		view.BgColor = w.colors.WindowBGColor.GetCUIAttr()
		view.FgColor = w.colors.FGColor.GetCUIAttr()
		view.SelBgColor = w.colors.WindowBGColor.GetCUIAttr()
		view.SetWritePos(0, w.h)
		view.Title = w.creature.Name
		view.Visible = false

	} else if err != nil {
		return err
	}

	return nil

}

func (w *ViewCreatureWidget) createKeybinds(g *gocui.Gui) error {

	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.onClose); err != nil {
		return err
	}
	return nil

}

func (w *ViewCreatureWidget) onClose(g *gocui.Gui, v *gocui.View) error {

	nextView, err := g.View(NameRootWidget)
	if err != nil {
		log.Println(err)
		return err
	}

	w.view.Visible = false
	w.view.Clear()

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(nextView.Name())
		return err
	})

	return nil
}
