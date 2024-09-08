package ui

import (
	"fmt"
	"slices"
	"strings"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

type SearchWidget struct {
	view          *gocui.View
	name          string
	title         string
	colors        *ColorPalette
	choices       map[string]string
	results       []string
	x, y          int
	w, h          int
	onSubmit      func(v string)
	killed        bool
	searchTerm    string
	selectedIndex int
}

func NewSearchWidget(name string, title string, colors *ColorPalette, choices map[string]string, onSubmit func(v string)) *SearchWidget {
	return &SearchWidget{
		name:          name,
		title:         title,
		choices:       choices,
		onSubmit:      onSubmit,
		colors:        colors,
		killed:        false,
		results:       []string{},
		searchTerm:    "",
		selectedIndex: 0,
	}
}

func (w *SearchWidget) Layout(g *gocui.Gui) error {
	if w.killed {
		g.DeleteKeybindings(w.name)
		g.DeleteView(w.name)
		w.view.Frame = false
		w.view.Clear()
		return nil
	}

	width, height := g.Size()

	w.w = utils.Clamp(width/3, height-5, width)
	w.h = height - 10

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		view.BgColor = w.colors.WindowBGColor.GetCUIAttr()
		view.FgColor = w.colors.FGColor.GetCUIAttr()
		view.Highlight = true
		view.Title = w.title
		view.SelBgColor = w.colors.FGColor.GetCUIAttr()
		view.SelFgColor = w.colors.BGColor.GetCUIAttr()
		w.setKeybinding(g)
	} else if err != nil {
		return err
	}

	w.results = []string{}

	for id, name := range w.choices {
		if strings.Contains(name, w.searchTerm) {
			w.results = append(w.results, id)
		}
	}

	view.Rewind()

	slices.Sort(w.results)

	for i, id := range w.results {

		if i > w.h/2 {
			break
		}
		w.view.SetWritePos(0, i+2)
		fmt.Fprintf(w.view, " %2d %s (%s)", i, w.choices[id], id)
	}

	view.SetCursor(0, w.selectedIndex+2)

	return nil

}

func (w *SearchWidget) setKeybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.Kill); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.move(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.move(-1)); err != nil {
		return err
	}

	return nil
}

func (w *SearchWidget) move(offset int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		w.selectedIndex = utils.Clamp(w.selectedIndex+offset, 0, w.h/2)
		w.Layout(g)
		return nil
	}
}

func (w *SearchWidget) Kill(g *gocui.Gui, v *gocui.View) error {
	w.killed = true

	g.Update(
		func(g *gocui.Gui) error {
			w.Layout(g)
			_, err := g.SetCurrentView(NameRootWidget)
			return err
		})
	return nil
}
