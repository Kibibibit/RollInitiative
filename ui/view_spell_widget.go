package ui

import (
	"fmt"
	"log"
	"strings"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

type ViewSpellWidget struct {
	view           *gocui.View
	name           string
	x, y           int
	w, h           int
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

	var spellLevelTitles = []string{
		"cantrip",
		"1st level",
		"2nd level",
		"3rd level",
		"4th level",
		"5th level",
		"6th level",
		"7th level",
		"8th level",
		"9th level",
	}

	var spellTraitNames = []string{
		"Casting Time",
		"Range",
		"Components",
		"Materials",
		"Duration",
	}

	width, height := g.Size()

	w.w = utils.Clamp(width-20, 70, width-4)
	w.h = height - 5

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

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
			w.spell.Name,
			w.colors.FgColor,
		),
	)

	view.SetWritePos(1, 2)

	fmt.Fprint(view, ApplyStyles(fmt.Sprintf(" - %s %s (%s)", spellLevelTitles[w.spell.Level], w.spell.School, w.spell.Source), gocui.AttrItalic))
	drawX, drawY := 1, 4

	if w.spell.Ritual {
		view.SetWritePos(1, 3)
		fmt.Fprint(view, ApplyStyles(" - Can be cast as ritual", gocui.AttrItalic))
		drawY += 1
	}

	spellTraits := []string{
		w.spell.CastingTime,
		w.spell.Range,
		w.spell.Components,
		w.spell.Materials,
		w.spell.Duration,
	}

	for i, trait := range spellTraits {
		if len(trait) > 0 {
			drawX, drawY = w.drawText(fmt.Sprintf("%s: %s", ApplyBold(spellTraitNames[i], w.colors.FgColor), trait), drawX, drawY)
		}

	}

	drawY += 1

	drawX, drawY = w.drawText(w.spell.Description, drawX, drawY)

	if len(w.spell.HigherLevels) > 0 {
		drawY += 1
		drawX, drawY = w.drawText(fmt.Sprintf("%s: %s", ApplyBold("At Higher Levels", w.colors.FgColor), w.spell.HigherLevels), drawX, drawY)
	}

	drawY += 1
	classesString := strings.Join(w.spell.Classes, ", ")

	drawX, drawY = w.drawText(fmt.Sprintf("%s: %s", ApplyBold("Spell Lists", w.colors.FgColor), classesString), drawX, drawY)

	return nil
}

func (w *ViewSpellWidget) drawText(text string, drawX, drawY int) (int, int) {
	return DrawText(w.view, w.w/2-4, w.h, w.colors, text, drawX, drawY)
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
