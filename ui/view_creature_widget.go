package ui

import (
	"fmt"
	"log"
	"strings"
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
		view.Title = "View Creature"
		view.Visible = true

	} else if err != nil {
		return err
	}

	//Draw Metadata

	view.SetWritePos(1, 1)

	fmt.Fprint(view, w.creature.Name)
	view.SetWritePos(1, 2)
	fmt.Fprintf(view, "- %s, %s, %s", w.creature.Type, w.creature.Size, w.creature.Alignment)

	view.SetWritePos(1, 4)
	fmt.Fprintf(view, "Armour Class: %s", w.creature.AC)

	conMod := (w.creature.Stats[2] - 10) / 2

	hpBoost := conMod * w.creature.HitDice

	view.SetWritePos(1, 5)
	hpString := fmt.Sprintf("Hit Points: %d (%dd%d + %d)", utils.AverageDiceRoll(w.creature.HitDice, w.creature.HitDiceType)+hpBoost, w.creature.HitDice, w.creature.HitDiceType, conMod*w.creature.HitDice)
	hpString = strings.ReplaceAll(hpString, " + -", " - ")
	hpString = strings.ReplaceAll(hpString, " + 0", "")
	fmt.Fprint(view, hpString)

	view.SetWritePos(1, 6)
	fmt.Fprintf(view, "Speed: %s", w.creature.Speed)

	//Draw statblock

	topLine := "┌── STR ──┬── DEX ──┬── CON ──┬── INT ──┬── WIS ──┬── CHA ──┐"
	bottomLine := "└─────────┴─────────┴─────────┴─────────┴─────────┴─────────┘"

	statLine := ""
	for _, value := range w.creature.Stats {
		statLine += fmt.Sprintf("│ %2d (%2d) ", value, ((value - 10) / 2))
	}
	statLine += "│"

	statLine = strings.ReplaceAll(statLine, "( ", "(+")
	view.SetWritePos(1, 8)
	fmt.Fprint(view, topLine)
	view.SetWritePos(1, 9)
	fmt.Fprint(view, statLine)
	view.SetWritePos(1, 10)
	fmt.Fprint(view, bottomLine)

	//Draw saves, resistences and so on

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
