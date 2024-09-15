package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type ConfirmWidget struct {
	view             *gocui.View
	name             string
	title            string
	colors           *ColorPalette
	x, y             int
	w, h             int
	onSubmit         func(bool)
	yesOption        string
	noOption         string
	currentSelection bool
	previousWidget   string
	killed           bool
	message          string
}

func NewConfirmWidget(title string, colors *ColorPalette, previousWidget string, currentSelection bool, message string, yesOption string, noOption string, onSubmit func(bool)) *ConfirmWidget {
	return &ConfirmWidget{
		name:             NameConfirmWidget,
		title:            title,
		colors:           colors,
		onSubmit:         onSubmit,
		yesOption:        yesOption,
		previousWidget:   previousWidget,
		noOption:         noOption,
		message:          message,
		killed:           false,
		currentSelection: currentSelection,
	}
}

func (w *ConfirmWidget) Layout(g *gocui.Gui) error {
	if w.killed {
		g.DeleteKeybindings(w.name)
		g.DeleteView(w.name)
		w.view.Clear()
		w.view.Frame = false

		return nil
	}

	width, height := g.Size()

	w.w = width / 4
	w.h = 3

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+1, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		view.BgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.Title = w.title
		view.SelBgColor = w.colors.FgColor.GetCUIAttr()
		view.SelFgColor = w.colors.BgColor.GetCUIAttr()
		w.setKeybinding(g)
	} else if err != nil {
		return err
	}

	w.view.Clear()
	w.view.Rewind()

	fmt.Fprint(w.view, w.message)

	view.SetWritePos(w.w/3, 2)

	yesDraw := w.yesOption
	noDraw := w.noOption
	if w.currentSelection {
		yesDraw = fmt.Sprintf("\x1b[7m%s\x1b[0m", w.yesOption)
	} else {
		noDraw = fmt.Sprintf("\x1b[7m%s\x1b[0m", w.noOption)
	}

	fmt.Fprint(view, yesDraw)

	view.SetWritePos((2*w.w)/3, 2)
	fmt.Fprint(view, noDraw)

	return nil
}

func (w *ConfirmWidget) setKeybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.Kill); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.onEnter); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.onKey); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.onKey); err != nil {
		return err
	}

	return nil
}

func (w *ConfirmWidget) Kill(g *gocui.Gui, v *gocui.View) error {
	w.killed = true

	g.Update(
		func(g *gocui.Gui) error {
			w.Layout(g)
			_, err := g.SetCurrentView(w.previousWidget)
			return err
		})
	return nil
}

func (w *ConfirmWidget) onKey(g *gocui.Gui, v *gocui.View) error {
	w.currentSelection = !w.currentSelection
	w.Layout(g)
	return nil
}

func (w *ConfirmWidget) onEnter(g *gocui.Gui, v *gocui.View) error {
	w.Kill(g, v)
	w.onSubmit(w.currentSelection)
	return nil
}
