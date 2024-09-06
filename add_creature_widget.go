package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type AddCreatureWidget struct {
	name              string
	w, h              int
	searchTerm        string
	filteredCreatures []string
}

func NewAddCreatureWidget(name string, w, h int) *AddCreatureWidget {
	return &AddCreatureWidget{name: name, w: w, h: h, searchTerm: "", filteredCreatures: []string{}}
}

type AddCreatureEditor struct {
	widget *AddCreatureWidget
}

func AddAddCreatureWidgetKeybinds(g *gocui.Gui) {
	if err := g.SetKeybinding(ADD_CREATURE_NAME, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if err := g.DeleteView(v.Name()); err != nil {
				return err
			}
			g.SetCurrentView(MAIN_TABLE_NAME)
			return nil
		}); err != nil {
		log.Panicln(err)
	}

}

func (e *AddCreatureEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		if ch < '0' || ch > '9' {
			e.widget.searchTerm += string(ch)
		}

	} else if key == gocui.KeySpace {
		e.widget.searchTerm += " "
	} else if key == gocui.KeyBackspace || key == gocui.KeyBackspace2 {
		if len(e.widget.searchTerm) > 0 {
			e.widget.searchTerm = e.widget.searchTerm[:len(e.widget.searchTerm)-1]
		}
	}

	e.widget.Search(v)

}

func (w *AddCreatureWidget) Layout(g *gocui.Gui) error {
	cols, lines := g.Size()

	x := int(cols / 2)
	y := int(lines / 2)

	width := int(w.w / 2)
	height := int(w.h / 2)

	x0 := int(x - width)
	x1 := int(x + width)
	y0 := int(y - height)
	y1 := int(y + height)

	view, err := g.SetView(w.name, x0, y0, x1, y1, 0)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Editable = true
		view.Editor = &AddCreatureEditor{widget: w}
		w.searchTerm = ""
		w.filteredCreatures = []string{}
		view.Title = "Add Creature"

	}

	w.Search(view)

	return nil

}

func (w *AddCreatureWidget) Search(v *gocui.View) {
	w.filteredCreatures = []string{}
	var index int = 0

	v.Clear()

	lowercaseSearchTerm := strings.ToLower(w.searchTerm)

	v.SetWritePos(3, 3)
	fmt.Fprint(v, "\x1b[37;2mName\x1b[0m")
	v.SetWritePos(w.w-4, 3)
	fmt.Fprint(v, "\x1b[37;2mCR\x1b[0m")

	for _, creatureId := range creatureIds {

		creature := creatureDict[creatureId]

		var contains bool = false

		if w.searchTerm == "" {
			contains = true
		} else if strings.Contains(strings.ToLower(creature.Name), lowercaseSearchTerm) {
			contains = true
		}

		if contains {
			w.filteredCreatures = append(w.filteredCreatures, creatureId)
			index += 1
			if index >= 10 {
				break
			}
		}
	}

	slices.Sort(w.filteredCreatures)

	for index, creatureId := range w.filteredCreatures {
		creature := creatureDict[creatureId]
		v.SetWritePos(1, index+4)
		drawIndex := index + 1
		if index+1 == 10 {
			drawIndex = 0
		}
		fmt.Fprintf(v, "\x1b[37;2m%d\x1b[0m %s", drawIndex, creature.Name)
		v.SetWritePos(w.w-4, index+4)
		fmt.Fprint(v, creature.CR)
		index += 1
		if index >= 10 {
			break
		}
	}

	v.SetWritePos(0, 0)
	v.SetLine(0, "")

	fmt.Fprintf(v, "🔎︎%s", w.searchTerm)
	v.SetLine(1, "")
	v.SetWritePos(0, 1)
	fmt.Fprintf(v, "\x1b[37;2mFound %d results\x1b[0m", len(w.filteredCreatures))
}

func newAddCreatureWindow(g *gocui.Gui, v *gocui.View, widget *AddCreatureWidget) error {

	widget.Layout(g)

	_, err := g.SetCurrentView(ADD_CREATURE_NAME)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	return nil
}
