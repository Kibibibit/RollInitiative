package ui

import (
	"fmt"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

// This widget will get a string

type StringInputWidget struct {
	view           *gocui.View
	name           string
	title          string
	colors         *ColorPalette
	x, y           int
	w, h           int
	onSubmit       func(v string)
	previousWidget string
	killed         bool
	data           string
}

func NewStringInputWidget(name string, title string, colors *ColorPalette, previousWidget string, onSubmit func(v string)) *StringInputWidget {
	return &StringInputWidget{
		name:           name,
		title:          title,
		onSubmit:       onSubmit,
		colors:         colors,
		killed:         false,
		data:           "",
		previousWidget: previousWidget,
	}
}

func (w *StringInputWidget) Layout(g *gocui.Gui) error {
	if w.killed {
		g.DeleteKeybindings(w.name)
		g.DeleteView(w.name)
		w.view.Clear()
		w.view.Frame = false

		return nil
	}

	width, height := g.Size()

	w.w = width / 4
	w.h = 1

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		view.BgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		// view.Highlight = true
		view.Title = w.title
		view.SelBgColor = w.colors.FgColor.GetCUIAttr()
		view.SelFgColor = w.colors.BgColor.GetCUIAttr()
		w.setKeybinding(g)
	} else if err != nil {
		return err
	}

	w.view.Clear()

	fmt.Fprint(w.view, w.data)

	return nil
}

func (w *StringInputWidget) setKeybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.Kill); err != nil {
		return err
	}

	for _, ch := range utils.ASCII_ALL {
		if err := g.SetKeybinding(w.name, ch, gocui.ModNone, w.onLetter(ch)); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.onLetter(' ')); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyBackspace|gocui.KeyBackspace2, gocui.ModNone, w.onBackspace); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.onEnter); err != nil {
		return err
	}

	return nil
}

func (w *StringInputWidget) Kill(g *gocui.Gui, v *gocui.View) error {
	w.killed = true

	g.Update(
		func(g *gocui.Gui) error {
			w.Layout(g)
			_, err := g.SetCurrentView(w.previousWidget)
			return err
		})
	return nil
}

func (w *StringInputWidget) onLetter(ch rune) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		w.data += string(ch)
		return w.Layout(g)
	}
}

func (w *StringInputWidget) onBackspace(g *gocui.Gui, v *gocui.View) error {
	if len(w.data) > 0 {
		w.data = w.data[0 : len(w.data)-1]
	}
	return w.Layout(g)
}

func (w *StringInputWidget) onEnter(g *gocui.Gui, v *gocui.View) error {
	w.onSubmit(w.data)
	return nil
}
