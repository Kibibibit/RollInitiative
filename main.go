package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
)

const TABLE_NAME = "table"
const ADD_CREATURE_NAME = "add_creature"

var (
	views     = []string{}
	beastiary *Beastiary
)

type MainTableWidget struct {
	name string
	x, y int
	data *IniativeTracker
}

func NewMainTableWidget(name string, x, y int) *MainTableWidget {

	return &MainTableWidget{
		name: name,
		x:    x,
		y:    y,
		data: &IniativeTracker{},
	}
}

func (w *MainTableWidget) Layout(g *gocui.Gui) error {

	width, height := g.Size()

	view, err := g.SetView(w.name, w.x, w.y, w.x+width-1, w.y+height-1, 0)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		g.SetCurrentView(w.name)

		for index, combatant := range w.data.combatants {
			view.SetWritePos(1, index)
			str := fmt.Sprintf("\x1b[37;2m%d\x1b[0m %s (%s)", index+1, combatant.creature.Name, combatant.tag)
			fmt.Fprintf(view, str)
		}

	}

	view.Title = "Roll Initiative"

	return nil
}

type AddCreatureWidget struct {
	name              string
	w, h              int
	beastiary         *Beastiary
	searchTerm        string
	filteredCreatures []Creature
}

func NewAddCreatureWidget(name string, w, h int) *AddCreatureWidget {
	return &AddCreatureWidget{name: name, w: w, h: h, beastiary: beastiary, searchTerm: "", filteredCreatures: []Creature{}}
}

type AddCreatureEditor struct {
	widget *AddCreatureWidget
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
	height := int(float32(lines)*0.8) / 2

	w.h = height

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
		w.filteredCreatures = []Creature{}
		view.Title = "Add Creature"
	}

	w.Search(view)

	return nil

}

func (w *AddCreatureWidget) Search(v *gocui.View) {
	w.filteredCreatures = []Creature{}
	var index int = 0

	v.Clear()

	lowercaseSearchTerm := strings.ToLower(w.searchTerm)

	for _, creature := range w.beastiary.Creatures {

		var contains bool = false

		if w.searchTerm == "" {
			contains = true
		} else if strings.Contains(strings.ToLower(creature.Name), lowercaseSearchTerm) {
			contains = true
		}

		if contains {
			w.filteredCreatures = append(w.filteredCreatures, creature)
			v.SetWritePos(1, index+3)
			fmt.Fprintf(v, "\x1b[37;2m%d\x1b[0m %s", index+1, creature.Name)
			index += 1
			if index >= w.h-5 {
				break
			}
		}
	}
	v.SetWritePos(0, 0)
	v.SetLine(0, "")

	fmt.Fprintf(v, "🔎︎%s", w.searchTerm)
	v.SetLine(1, "")
	v.SetWritePos(0, 1)
	fmt.Fprintf(v, "\x1b[37;2mFound %d results\x1b[0m", len(w.filteredCreatures))
}

func main() {

	b, err := LoadBeastiary("./data/creatures.xml")
	beastiary = b

	if err != nil {
		log.Panicln(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	table := NewMainTableWidget(TABLE_NAME, 0, 0)

	addCreature := NewAddCreatureWidget(ADD_CREATURE_NAME, 50, 0)

	g.SetManager(table)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(TABLE_NAME, 'a', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newAddCreatureWindow(g, v, addCreature)
		}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(ADD_CREATURE_NAME, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if err := g.DeleteView(v.Name()); err != nil {
				return err
			}
			g.SetCurrentView(TABLE_NAME)
			return nil
		}); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
