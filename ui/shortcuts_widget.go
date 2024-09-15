package ui

import (
	"fmt"
	"log"
	"strconv"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

const (
	shortcutsWidgetHeight  int  = 5
	shortcutsWidgetNilMenu rune = '_'
	shortcutWiki           rune = 'w'
	shortcutEdit           rune = 'e'
	shortcutAdd            rune = 'a'
	shortcutTurn           rune = 't'
)

type ShortcutsWidget struct {
	view         *gocui.View
	name         string
	x, y         int
	w, h         int
	dataStore    *models.DataStore
	colors       *ColorPalette
	submenu      rune
	rootWidget   *RootWidget
	shortcuts    map[rune]map[rune]*Shortcut
	submenuNames map[rune]string
}

type Shortcut struct {
	name         string
	onPress      func(g *gocui.Gui, v *gocui.View) error
	entryOnly    bool
	creatureOnly bool
}

func NewShortcutsWidget(rootWidget *RootWidget, dataStore *models.DataStore, colors *ColorPalette) *ShortcutsWidget {
	out := ShortcutsWidget{
		rootWidget: rootWidget,
		name:       NameShortcutsWidget,
		dataStore:  dataStore,
		submenu:    shortcutsWidgetNilMenu,
		colors:     colors,
	}

	submenuNamesDict := map[rune]string{
		shortcutAdd:  "Add",
		shortcutEdit: "Edit",
		shortcutWiki: "Wiki",
		shortcutTurn: "Turn",
	}

	shortcutsDict := map[rune]map[rune]*Shortcut{
		shortcutAdd: {
			'c': {"Creature", out.addCreatureEntry, false, false},
			'p': {"Party", out.addPartyEntries, false, false},
		},
		shortcutEdit: {
			'r': {"Remove", out.deleteCreatureEntry, true, false},
			'h': {"Health", out.editCreatureHealth, true, true},
			'd': {"Damage/Heal", out.damageCreature, true, true},
			's': {"Status", out.editCreatureStatus, true, false},
			'i': {"Initiative", out.editCreatureIniative, true, false},
		},
		shortcutWiki: {
			'c': {"Creatures", out.openCreatureWiki, false, false},
			's': {"Spells", out.openSpellsWiki, false, false},
			'e': {"Current", out.openCurrentWiki, true, true},
		},
		shortcutTurn: {
			'n': {"Next", out.moveTurn(1), true, false},
			'p': {"Previous", out.moveTurn(-1), true, false},
			's': {"Set", out.setTurn, true, false},
		},
	}

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

	colW := w.w / 5

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		view.Visible = false
		w.createKeybinds(g)
		view.Frame = false
		view.BgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.SetWritePos(0, 0)

	} else if err != nil {
		return err
	}
	view.Clear()

	view.Rewind()

	items := []string{}

	drawX, drawY := 1, 1

	if w.submenu == '_' {

		for key, name := range w.submenuNames {
			items = append(items, fmt.Sprintf("%c %s", key, name))
		}

	} else {
		for key, shortcut := range w.shortcuts[w.submenu] {

			dimLine := (!w.rootWidget.CurrentEntryIsNotPlayer() && shortcut.creatureOnly) ||
				(!w.rootWidget.ValidCurrentEntry() && shortcut.entryOnly)
			drawLine := fmt.Sprintf("%c %s", key, shortcut.name)

			if dimLine {
				drawLine = ApplyFgColor(drawLine, w.colors.FgColorDim)
			}
			items = append(items, drawLine)
		}
	}

	for _, item := range items {
		view.SetWritePos(drawX, drawY)
		fmt.Fprintf(view, "%s\n", item)
		drawY += 2
		if drawY >= w.h {
			drawY = 1
			drawX += colW
		}
	}

	return nil

}

func (w *ShortcutsWidget) createKeybinds(g *gocui.Gui) error {

	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.onClose); err != nil {
		return err
	}

	for _, ch := range utils.ASCII_LETTERS {
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
				w.view.Clear()
			}
		} else {
			if shortcut, ok := w.shortcuts[w.submenu][key]; ok {
				if shortcut.creatureOnly && !w.rootWidget.CurrentEntryIsNotPlayer() {
					return w.badShortcut(g, v)
				}
				if shortcut.entryOnly && !w.rootWidget.ValidCurrentEntry() {
					return w.badShortcut(g, v)
				}

				return shortcut.onPress(g, v)

			}
		}

		return nil
	}
}

func (w *ShortcutsWidget) openCreatureWiki(g *gocui.Gui, v *gocui.View) error {

	w.hide()

	NewCreatureSearch(g, w.colors, w.dataStore, func(result string) {

		viewCreature := NewViewCreatureWidget(w.dataStore, NameRootWidget, w.colors, w.dataStore.GetCreature(result))

		viewCreature.Layout(g)

		g.Update(func(g *gocui.Gui) error {

			_, err := g.SetCurrentView(viewCreature.name)
			return err
		})

	})

	return nil
}

func (w *ShortcutsWidget) openSpellsWiki(g *gocui.Gui, v *gocui.View) error {

	w.hide()

	NewSpellSearch(g, w.colors, w.dataStore, func(result string) {

		viewSpell := NewViewSpellWidget(w.dataStore, NameRootWidget, w.colors, w.dataStore.GetSpell(result))

		viewSpell.Layout(g)

		g.Update(func(g *gocui.Gui) error {

			_, err := g.SetCurrentView(viewSpell.name)
			return err
		})

	})

	return nil
}

func (w *ShortcutsWidget) openCurrentWiki(g *gocui.Gui, v *gocui.View) error {
	w.hide()

	initEntry := w.dataStore.IniativeEntries[w.rootWidget.GetCurrentEntryId()]
	creature := w.dataStore.GetCreature(initEntry.CreatureId)
	if creature == nil {
		return nil
	}
	viewCreature := NewViewCreatureWidget(w.dataStore, NameRootWidget, w.colors, creature)

	viewCreature.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(viewCreature.name)
		return err
	})

	return nil

}

func (w *ShortcutsWidget) addCreatureEntry(g *gocui.Gui, v *gocui.View) error {
	w.hide()

	NewCreatureSearch(g, w.colors, w.dataStore, func(result string) {

		creature := w.dataStore.GetCreature(result)
		var addCreatureWidget *AddCreatureWidget
		addCreatureWidget = NewAddCreatureWidget(
			w.colors,
			creature.Name,
			func(rollHp bool, count int, tags []string) {
				addCreatureWidget.Kill(g, v)

				for _, tag := range tags {
					w.dataStore.NewCreatureEntry(result, tag, rollHp)
				}

				w.rootWidget.Layout(g)

				g.Update(func(g *gocui.Gui) error {

					_, err := g.SetCurrentView(NameRootWidget)
					return err
				})
			},
		)

		addCreatureWidget.Layout(g)
		g.Update(func(g *gocui.Gui) error {

			_, err := g.SetCurrentView(addCreatureWidget.name)
			return err
		})

	})
	return nil
}

func (w *ShortcutsWidget) deleteCreatureEntry(g *gocui.Gui, v *gocui.View) error {
	w.hide()

	w.dataStore.DeleteCreatureEntry(w.rootWidget.GetCurrentEntryId())

	w.rootWidget.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(NameRootWidget)
		return err
	})

	return nil
}

func (w *ShortcutsWidget) editCreatureHealth(g *gocui.Gui, v *gocui.View) error {
	w.hide()

	entry := w.dataStore.IniativeEntries[w.rootWidget.GetCurrentEntryId()]

	stringWidget := NewStringInputWidget(
		NameStringWidget,
		"Edit creature health",
		w.colors,
		NameRootWidget,
		utils.ASCII_NUMBERS+"-",
		fmt.Sprintf("%d", entry.Hp),
		func(result string) {
			if len(result) > 0 {

				i64, err := strconv.ParseInt(result, 10, 64)
				if err != nil {
					log.Panicln(err)
				}
				entry.Hp = int(i64)

			}
		},
	)

	stringWidget.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(stringWidget.name)
		return err
	})

	return nil
}

func (w *ShortcutsWidget) editCreatureIniative(g *gocui.Gui, v *gocui.View) error {
	w.hide()

	entry := w.dataStore.IniativeEntries[w.rootWidget.GetCurrentEntryId()]

	stringWidget := NewStringInputWidget(
		NameStringWidget,
		"Edit creature initiative",
		w.colors,
		NameRootWidget,
		utils.ASCII_NUMBERS+"-",
		fmt.Sprintf("%d", entry.IniativeRoll),
		func(result string) {
			if len(result) > 0 {

				i64, err := strconv.ParseInt(result, 10, 64)
				if err != nil {
					log.Panicln(err)
				}
				entry.IniativeRoll = int(i64)

			}
		},
	)

	stringWidget.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(stringWidget.name)
		return err
	})

	return nil
}

func (w *ShortcutsWidget) damageCreature(g *gocui.Gui, v *gocui.View) error {
	w.hide()

	entry := w.dataStore.IniativeEntries[w.rootWidget.GetCurrentEntryId()]

	stringWidget := NewStringInputWidget(
		NameStringWidget,
		"Damage Creature",
		w.colors,
		NameRootWidget,
		utils.ASCII_NUMBERS+"-",
		"",
		func(result string) {
			if len(result) > 0 {

				i64, err := strconv.ParseInt(result, 10, 64)
				if err != nil {
					log.Panicln(err)
				}
				entry.Hp = utils.Clamp(entry.Hp-int(i64), -entry.MaxHp, entry.MaxHp)

			}
		},
	)

	stringWidget.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(stringWidget.name)
		return err
	})

	return nil
}

func (w *ShortcutsWidget) editCreatureStatus(g *gocui.Gui, v *gocui.View) error {
	w.hide()

	entry := w.dataStore.IniativeEntries[w.rootWidget.GetCurrentEntryId()]

	stringWidget := NewStringInputWidget(NameStringWidget, "Edit creature status", w.colors, NameRootWidget, utils.ASCII_ALL, entry.Statuses, func(result string) {
		entry.Statuses = result
	})

	stringWidget.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(stringWidget.name)
		return err
	})

	return nil
}

func (w *ShortcutsWidget) addPartyEntries(g *gocui.Gui, v *gocui.View) error {
	w.hide()

	NewPartySearch(g, w.colors, w.dataStore, func(result string) {

		w.dataStore.NewPartyEntries(result)

	})
	return nil
}

func (w *ShortcutsWidget) badShortcut(g *gocui.Gui, _ *gocui.View) error {
	w.hide()

	w.rootWidget.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(NameRootWidget)
		return err
	})

	return nil
}

func (w *ShortcutsWidget) moveTurn(offset int) func(g *gocui.Gui, _ *gocui.View) error {

	return func(g *gocui.Gui, _ *gocui.View) error {
		w.hide()
		w.rootWidget.currentTurnIndex += offset
		if w.rootWidget.currentTurnIndex < 0 {
			w.rootWidget.currentTurnIndex += len(w.rootWidget.entryIds)
		}
		if w.rootWidget.currentTurnIndex >= len(w.rootWidget.entryIds) {
			w.rootWidget.currentTurnIndex -= len(w.rootWidget.entryIds)
		}

		w.rootWidget.Layout(g)

		g.Update(func(g *gocui.Gui) error {

			_, err := g.SetCurrentView(NameRootWidget)
			return err
		})

		return nil
	}

}

func (w *ShortcutsWidget) setTurn(g *gocui.Gui, _ *gocui.View) error {

	w.hide()

	w.rootWidget.currentTurnIndex = w.rootWidget.currentEntryIndex

	w.rootWidget.Layout(g)

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(NameRootWidget)
		return err
	})

	return nil
}
