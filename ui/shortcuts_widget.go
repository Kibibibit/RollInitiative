package ui

import (
	"fmt"
	"windmills/roll_initiative/models"

	"github.com/awesome-gocui/gocui"
)

const (
	shortcutsWidgetHeight  int  = 5
	shortcutsWidgetNilMenu rune = '_'
)

type ShortcutsWidget struct {
	view         *gocui.View
	name         string
	x, y         int
	w, h         int
	dataStore    *models.DataStore
	colors       *ColorPalette
	submenu      rune
	shortcuts    map[rune]map[rune]*Shortcut
	submenuNames map[rune]string
}

type Shortcut struct {
	name    string
	onPress func(g *gocui.Gui, v *gocui.View) error
}

func NewShortcutsWidget(dataStore *models.DataStore, colors *ColorPalette) *ShortcutsWidget {
	out := ShortcutsWidget{
		name:      NameShortcutsWidget,
		dataStore: dataStore,
		submenu:   shortcutsWidgetNilMenu,
		colors:    colors,
	}

	submenuNamesDict := make(map[rune]string)
	submenuNamesDict['w'] = "Wiki"

	shortcutsDict := make(map[rune]map[rune]*Shortcut)

	shortcutsWikiDict := make(map[rune]*Shortcut)

	shortcutsWikiDict['c'] = &Shortcut{"Creatures", out.openCreatureWiki}
	shortcutsWikiDict['s'] = &Shortcut{"Spells", out.openSpellsWiki}

	shortcutsDict['w'] = shortcutsWikiDict

	out.submenuNames = submenuNamesDict
	out.shortcuts = shortcutsDict

	return &out
}

func (w *ShortcutsWidget) Layout(g *gocui.Gui) error {

	width, height := g.Size()

	w.x = -1
	w.y = height - 1 - shortcutsWidgetHeight

	w.w = width - 1
	w.h = shortcutsWidgetHeight - 1

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		view.Visible = false
		w.createKeybinds(g)
		view.Frame = false
		view.BgColor = w.colors.WindowBGColor.GetCUIAttr()
		view.FgColor = w.colors.FGColor.GetCUIAttr()
		view.SetWritePos(0, 0)

	} else if err != nil {
		return err
	}

	view.Rewind()

	items := []string{}

	if w.submenu == '_' {

		for key, name := range w.submenuNames {
			items = append(items, fmt.Sprintf("%c %s", key, name))
		}

	} else {
		for key, shortcut := range w.shortcuts[w.submenu] {
			items = append(items, fmt.Sprintf("%c %s", key, shortcut.name))
		}
	}

	for _, item := range items {
		fmt.Fprintf(view, "%s\n", item)
	}

	return nil

}

func (w *ShortcutsWidget) createKeybinds(g *gocui.Gui) error {

	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.onClose); err != nil {
		return err
	}

	for ch := 'A'; ch <= 'z'; ch++ {
		if err := g.SetKeybinding(w.name, ch, gocui.ModNone, w.onKeypress(ch)); err != nil {
			return err
		}
	}

	return nil

}

func (w *ShortcutsWidget) hide() {
	w.view.Visible = false
	w.submenu = shortcutsWidgetNilMenu
	w.view.Clear()
}

func (w *ShortcutsWidget) onClose(g *gocui.Gui, v *gocui.View) error {
	rootView, err := g.View(NameRootWidget)

	if err != nil {
		return err
	}

	w.hide()

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(rootView.Name())
		return err
	})

	return nil
}

func (w *ShortcutsWidget) onKeypress(key rune) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		if w.submenu == shortcutsWidgetNilMenu {
			if _, ok := w.submenuNames[key]; ok {
				w.submenu = key
			}
		} else {
			if shortcut, ok := w.shortcuts[w.submenu][key]; ok {
				return shortcut.onPress(g, v)
			}
		}

		return nil
	}
}

func (w *ShortcutsWidget) openCreatureWiki(g *gocui.Gui, v *gocui.View) error {

	w.hide()
	//TODO: Put this into a function in search widget
	var creatureSearch *SearchWidget

	creatureSearch = NewSearchWidget(
		NameSearchCreaturesWidget,
		"Find Creature",
		w.colors,
		w.dataStore.CreatureNames,
		func(result string) {
			creatureSearch.Kill(g, v)
		},
	)

	creatureSearch.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(creatureSearch.name)
		return err
	})

	return nil
}

func (w *ShortcutsWidget) openSpellsWiki(g *gocui.Gui, v *gocui.View) error {
	return w.onClose(g, v)
}
