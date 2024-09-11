package ui

import (
	"fmt"
	"slices"
	"strings"
	"windmills/roll_initiative/models"

	"github.com/awesome-gocui/gocui"
)

const (
	keybindBufferLength = 2
	keybindLeader       = gocui.KeySpace
)

type RootWidget struct {
	view          *gocui.View
	name          string
	x, y          int
	w, h          int
	dataStore     *models.DataStore
	colors        *ColorPalette
	keybindBuffer string
	leaderPressed bool
}

func NewRootWidget(dataStore *models.DataStore, colors *ColorPalette) *RootWidget {
	return &RootWidget{
		name:          NameRootWidget,
		dataStore:     dataStore,
		x:             -1,
		y:             -1,
		keybindBuffer: "",
		leaderPressed: false,
		colors:        colors,
	}
}

func (w *RootWidget) Layout(g *gocui.Gui) error {

	var columnLengths = []int{50, 15, 5, 20}

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

	} else if err != nil {
		return err
	}

	table := [][]string{}

	table = append(table, []string{"Combatant", "Initiative", "HP", "Statuses"})

	entries := []*models.IniativeEntry{}

	for _, x := range w.dataStore.IniativeEntries {
		entries = append(entries, x)
	}

	slices.SortFunc(entries, func(a *models.IniativeEntry, b *models.IniativeEntry) int {
		return b.IniativeRoll - a.IniativeRoll
	})

	for _, entry := range entries {
		creature := w.dataStore.GetCreature(entry.CreatureId)
		var row []string

		if entry.IsPlayer {
			row = []string{entry.CreatureId, fmt.Sprintf("%d", entry.IniativeRoll), "", entry.Statuses}
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
		view.SetWritePos(1, y)
		fmt.Fprint(view, line)
	}

	return nil

}

func (w *RootWidget) createKeybinds(g *gocui.Gui) error {

	if err := g.SetKeybinding(w.name, keybindLeader, gocui.ModNone, w.onLeader); err != nil {
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
