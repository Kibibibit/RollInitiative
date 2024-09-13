package ui

import (
	"fmt"
	"slices"
	"strings"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

const (
	keybindBufferLength = 2
	keybindLeader       = gocui.KeySpace
)

type RootWidget struct {
	view              *gocui.View
	name              string
	x, y              int
	w, h              int
	dataStore         *models.DataStore
	colors            *ColorPalette
	keybindBuffer     string
	leaderPressed     bool
	entryIds          []int
	currentEntryIndex int
}

func NewRootWidget(dataStore *models.DataStore, colors *ColorPalette) *RootWidget {
	return &RootWidget{
		name:              NameRootWidget,
		dataStore:         dataStore,
		x:                 -1,
		y:                 -1,
		keybindBuffer:     "",
		leaderPressed:     false,
		colors:            colors,
		entryIds:          []int{},
		currentEntryIndex: 0,
	}
}

func (w *RootWidget) Layout(g *gocui.Gui) error {

	width, height := g.Size()

	w.w = width
	w.h = height

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		w.createKeybinds(g)
		view.Frame = false
		view.BgColor = w.colors.BgColor.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.SelFgColor = w.colors.BgColor.GetCUIAttr()
		view.SelBgColor = w.colors.FgColor.GetCUIAttr()

	} else if err != nil {
		return err
	}

	view.Clear()

	table := [][]string{}

	table = append(table, []string{"Combatant", "Initiative", "HP", "Statuses"})

	w.entryIds = []int{}

	for entryId := range w.dataStore.IniativeEntries {
		w.entryIds = append(w.entryIds, entryId)
	}

	slices.SortStableFunc(w.entryIds, func(a int, b int) int {

		cA := w.dataStore.IniativeEntries[a]
		cB := w.dataStore.IniativeEntries[b]

		if cA.IniativeRoll != cB.IniativeRoll {
			return cB.IniativeRoll - cA.IniativeRoll
		} else if cB.DexScore != cA.DexScore {
			return cB.DexScore - cA.DexScore
		} else {
			return strings.Compare(fmt.Sprintf("%s %s", cA.CreatureId, cA.Tag), fmt.Sprintf("%s %s", cB.CreatureId, cB.Tag))
		}
	})

	for _, entryId := range w.entryIds {
		entry := w.dataStore.IniativeEntries[entryId]
		creature := w.dataStore.GetCreature(entry.CreatureId)
		var row []string

		if entry.IsPlayer {
			player := w.dataStore.GetPlayer(entry.CreatureId)
			row = []string{player.Name, fmt.Sprintf("%d", entry.IniativeRoll), "", entry.Statuses}
		} else {

			name := creature.Name
			if len(entry.Tag) > 0 {
				name = fmt.Sprintf("%s (%s)", name, entry.Tag)
			}

			row = []string{
				name,
				fmt.Sprintf("%d", entry.IniativeRoll),
				fmt.Sprintf("%d", entry.Hp),
				entry.Statuses,
			}
		}

		table = append(table, row)

	}

	//Draw table
	drawLines := []string{}

	var midLine string
	var endLine string

	borderLines := []string{}
	columnLengths := []int{}

	for _, row := range table {
		for x, cell := range row {
			if len(columnLengths)-1 < x {
				columnLengths = append(columnLengths, 0)
			}

			if len(cell) > columnLengths[x] {
				columnLengths[x] = len(cell)
			}
		}
	}

	for _, length := range columnLengths {
		borderLines = append(borderLines, strings.Repeat("─", length+2))

	}
	drawLines = append(drawLines, fmt.Sprintf("┌%s┐", strings.Join(borderLines, "┬")))
	midLine = fmt.Sprintf("├%s┤", strings.Join(borderLines, "┼"))
	endLine = fmt.Sprintf("└%s┘", strings.Join(borderLines, "┴"))

	for y, row := range table {
		rowCells := []string{}
		for x, cell := range row {
			for len(cell) < columnLengths[x] {
				cell = fmt.Sprintf("%s ", cell)
			}
			if y == 0 {
				cell = ApplyBold(cell, w.colors.FgColor)
			}
			rowCells = append(rowCells, fmt.Sprintf(" %s ", cell))
		}
		drawLines = append(drawLines, fmt.Sprintf("│%s│", strings.Join(rowCells, "│")))
		if y == 0 {
			drawLines = append(drawLines, midLine)
		}
	}

	drawLines = append(drawLines, endLine)

	for y, line := range drawLines {
		view.SetWritePos(0, y)
		fmt.Fprint(view, line)
	}

	if len(w.dataStore.IniativeEntries) > 0 {
		w.currentEntryIndex = utils.Clamp(w.currentEntryIndex, 0, len(w.dataStore.IniativeEntries)-1)

		view.Highlight = true
		view.SetCursor(0, w.currentEntryIndex+3)
	} else {
		view.Highlight = false
		w.currentEntryIndex = -1
	}

	return nil

}

func (w *RootWidget) createKeybinds(g *gocui.Gui) error {

	if err := g.SetKeybinding(w.name, keybindLeader, gocui.ModNone, w.onLeader); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.moveCursor(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.moveCursor(1)); err != nil {
		return err
	}

	return nil

}

func (w *RootWidget) onLeader(g *gocui.Gui, v *gocui.View) error {
	shortcutView, err := g.View(NameShortcutsWidget)

	if err != nil {
		return err
	}

	shortcutView.Visible = true

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(shortcutView.Name())
		return err
	})

	return nil
}

func (w *RootWidget) GetCurrentEntryId() int {
	if len(w.entryIds) > 0 {
		return w.entryIds[w.currentEntryIndex]
	}
	return -1
}

func (w *RootWidget) CurrentEntryIsNotPlayer() bool {
	if len(w.entryIds) > 0 {
		creatureId := w.dataStore.IniativeEntries[w.GetCurrentEntryId()].CreatureId
		return w.dataStore.GetCreature(creatureId) != nil
	}
	return false
}

func (w *RootWidget) ValidCurrentEntry() bool {
	return len(w.entryIds) > 0
}

func (w *RootWidget) moveCursor(offset int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		w.currentEntryIndex += offset
		if w.currentEntryIndex < 0 {
			w.currentEntryIndex += len(w.dataStore.IniativeEntries)
		}
		if w.currentEntryIndex >= len(w.dataStore.IniativeEntries) {
			w.currentEntryIndex -= len(w.dataStore.IniativeEntries)
		}
		w.Layout(g)
		return nil

	}
}
