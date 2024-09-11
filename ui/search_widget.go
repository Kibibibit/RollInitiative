package ui

import (
	"fmt"
	"slices"
	"strings"
	"windmills/roll_initiative/models"
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
	renderItem    func(int, int, string) string
	killed        bool
	searchTerm    string
	selectedIndex int
}

func NewSearchWidget(name string, title string, colors *ColorPalette, choices map[string]string, renderItem func(int, int, string) string, onSubmit func(v string)) *SearchWidget {
	return &SearchWidget{
		name:          name,
		title:         title,
		choices:       choices,
		onSubmit:      onSubmit,
		renderItem:    renderItem,
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
		view.BgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.Highlight = true
		view.Title = w.title
		view.SelBgColor = w.colors.FgColor.GetCUIAttr()
		view.SelFgColor = w.colors.BgColor.GetCUIAttr()
		w.setKeybinding(g)
	} else if err != nil {
		return err
	}

	w.view.Clear()
	w.results = []string{}

	for id, name := range w.choices {
		if strings.Contains(strings.ToLower(name), strings.ToLower(w.searchTerm)) {
			w.results = append(w.results, id)
		}
	}

	w.selectedIndex = utils.Clamp(w.selectedIndex, 0, len(w.results)-1)

	w.view.Highlight = len(w.results) != 0

	view.Rewind()

	fmt.Fprintf(w.view, "Search: \x1b[4m%s\x1b[0m", w.searchTerm)

	slices.Sort(w.results)

	for i, id := range w.results {

		if i > w.h-2 {
			break
		}
		w.view.SetWritePos(0, i+2)
		fmt.Fprint(w.view, w.renderItem(i, w.w, id))
	}
	view.SetCursor(0, w.selectedIndex+2)

	return nil

}

func (w *SearchWidget) setKeybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.Kill); err != nil {
		return err
	}

	for ch := 'A'; ch <= 'z'; ch++ {
		if err := g.SetKeybinding(w.name, ch, gocui.ModNone, w.keypress(ch)); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.keypress(' ')); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyBackspace|gocui.KeyBackspace2, gocui.ModNone, w.backspace); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.move(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.move(-1)); err != nil {
		return err
	}

	for ch := '0'; ch <= '9'; ch++ {
		if err := g.SetKeybinding(w.name, ch, gocui.ModNone, w.onSelect(int(ch)-0x30)); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.onEnter); err != nil {
		return err
	}

	return nil
}

func (w *SearchWidget) move(offset int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		w.selectedIndex = utils.Clamp(w.selectedIndex+offset, 0, w.h-2)
		w.Layout(g)
		return nil
	}
}

func (w *SearchWidget) keypress(ch rune) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		w.searchTerm += string(ch)
		w.Layout(g)
		return nil
	}
}

func (w *SearchWidget) onSelect(index int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		w.onSubmit(w.results[index])

		return nil
	}
}

func (w *SearchWidget) onEnter(g *gocui.Gui, v *gocui.View) error {
	return w.onSelect(w.selectedIndex)(g, v)
}

func (w *SearchWidget) backspace(g *gocui.Gui, v *gocui.View) error {
	if len(w.searchTerm) > 0 {
		w.searchTerm = w.searchTerm[0 : len(w.searchTerm)-1]
	}
	w.Layout(g)
	return nil
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

func NewCreatureSearch(g *gocui.Gui, colors *ColorPalette, dataStore *models.DataStore, onSubmit func(string)) {
	var creatureSearch *SearchWidget

	creatureSearch = NewSearchWidget(
		NameSearchCreaturesWidget,
		"Find Creature",
		colors,
		dataStore.CreatureNames,
		func(index int, width int, id string) string {
			c := dataStore.GetCreature(id)
			return RenderCreatureSearchRow(c, colors, index, width)
		},
		func(result string) {
			creatureSearch.Kill(g, creatureSearch.view)
			onSubmit(result)
		},
	)

	creatureSearch.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(creatureSearch.name)
		return err
	})
}

func NewSpellSearch(g *gocui.Gui, colors *ColorPalette, dataStore *models.DataStore, onSubmit func(string)) {
	var spellSearch *SearchWidget

	spellSearch = NewSearchWidget(
		NameSearchSpellsWidget,
		"Find Spell",
		colors,
		dataStore.SpellNames,
		func(index int, width int, id string) string {
			s := dataStore.GetSpell(id)
			return RenderSpellSearchRow(s, colors, index, width)
		},
		func(result string) {
			spellSearch.Kill(g, spellSearch.view)
			onSubmit(result)
		},
	)

	spellSearch.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(spellSearch.name)
		return err
	})
}
