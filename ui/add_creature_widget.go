package ui

import (
	"fmt"
	"log"
	"strings"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

type AddCreatureWidget struct {
	view          *gocui.View
	name          string
	colors        *ColorPalette
	creatureName  string
	x, y          int
	w, h          int
	onSubmit      func(bool, int, []string)
	rollHp        bool
	tags          []string
	killed        bool
	selectedField int
	ROLL_HP_FIELD int
	TAGS_FIELD    int
	DONE_FIELD    int
}

func NewAddCreatureWidget(colors *ColorPalette, creatureName string, onSubmit func(bool, int, []string)) *AddCreatureWidget {
	return &AddCreatureWidget{
		name:          NameAddCreatureWidget,
		colors:        colors,
		creatureName:  creatureName,
		onSubmit:      onSubmit,
		rollHp:        true,
		tags:          []string{},
		killed:        false,
		ROLL_HP_FIELD: 0,
		TAGS_FIELD:    1,
		DONE_FIELD:    2,
	}
}

func (w *AddCreatureWidget) Layout(g *gocui.Gui) error {

	if w.killed {
		g.DeleteKeybindings(w.name)
		g.DeleteView(w.name)
		w.view.Frame = false
		w.view.Clear()
		return nil
	}

	width, height := g.Size()

	w.w = utils.Clamp(width/3, height-5, width)
	w.h = 10

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		view.BgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.Highlight = true
		view.Title = "Add Creatures"
		view.SelBgColor = w.colors.FgColor.GetCUIAttr()
		view.SelFgColor = w.colors.BgColor.GetCUIAttr()
		view.Wrap = true

		if err2 := w.setKeybinding(g); err2 != nil {
			log.Panicln(err2)
		}
	} else if err != nil {
		return err
	}

	w.view.Clear()

	w.view.Rewind()

	w.view.SetWritePos(1, 0)

	fmt.Fprintf(view, "Adding new: %s", w.creatureName)

	//Draw roll hp field
	w.view.SetWritePos(1, 3)

	rollHpBoldColor := w.colors.FgColor
	if w.selectedField == w.ROLL_HP_FIELD {
		rollHpBoldColor = w.colors.BgColorWindow
	}
	trueFalseString := "Yes"
	if w.rollHp == false {
		trueFalseString = "No"
	}
	rollHpString := fmt.Sprintf("%s %s", ApplyBold("Roll Creature HP:", rollHpBoldColor), trueFalseString)

	fmt.Fprint(view, rollHpString)

	//Draw tag field
	w.view.SetWritePos(1, 4)

	tagBoldColor := w.colors.FgColor
	if w.selectedField == w.TAGS_FIELD {
		tagBoldColor = w.colors.BgColorWindow
	}
	tagString := fmt.Sprintf("%s %s", ApplyBold("Creature Tags:", tagBoldColor), strings.Join(w.tags, ", "))

	fmt.Fprint(view, tagString)

	//Draw done button
	w.view.SetWritePos(1, 6)

	doneBoldColor := w.colors.FgColor
	if w.selectedField == w.DONE_FIELD {
		doneBoldColor = w.colors.BgColorWindow
	}
	doneString := ApplyBold("DONE", doneBoldColor)

	fmt.Fprint(view, doneString)

	if w.selectedField == w.ROLL_HP_FIELD {
		view.SetCursor(0, 3)
	} else if w.selectedField == w.TAGS_FIELD {
		view.SetCursor(0, 4)
	} else if w.selectedField == w.DONE_FIELD {
		view.SetCursor(0, 6)
	}

	return nil
}

func (w *AddCreatureWidget) setKeybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.Kill); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.moveField(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.moveField(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.onEnter); err != nil {
		return err
	}

	return nil
}

func (w *AddCreatureWidget) moveField(offset int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		w.selectedField = utils.Clamp(w.selectedField+offset, w.ROLL_HP_FIELD, w.DONE_FIELD)
		w.Layout(g)
		return nil
	}
}

func (w *AddCreatureWidget) onEnter(g *gocui.Gui, v *gocui.View) error {

	if w.selectedField == w.DONE_FIELD {

		if len(w.tags) < 1 {
			w.tags = []string{""}
		}

		w.onSubmit(w.rollHp, len(w.tags), w.tags)
	} else if w.selectedField == w.ROLL_HP_FIELD {

		w.rollHp = !w.rollHp
		w.Layout(g)

	} else if w.selectedField == w.TAGS_FIELD {

		w.tags = []string{}

		w.getStrings(g)

	}

	return nil
}

func (w *AddCreatureWidget) getStrings(g *gocui.Gui) {

	var inputWidget *StringInputWidget

	inputWidget = NewStringInputWidget(
		NameStringWidget, "Input Tags", w.colors, w.name, func(result string) {

			w.tags = strings.Split(result, ",")

			inputWidget.killed = true

			inputWidget.Layout(g)

			w.Layout(g)
			g.Update(
				func(g *gocui.Gui) error {

					_, err := g.SetCurrentView(w.name)
					return err
				})

		})

	inputWidget.Layout(g)

	g.Update(
		func(g *gocui.Gui) error {

			_, err := g.SetCurrentView(inputWidget.name)
			return err
		})

}

func (w *AddCreatureWidget) Kill(g *gocui.Gui, v *gocui.View) error {
	w.killed = true

	g.Update(
		func(g *gocui.Gui) error {
			w.Layout(g)
			_, err := g.SetCurrentView(NameRootWidget)
			return err
		})
	return nil
}
