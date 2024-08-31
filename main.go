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
	name       string
	w, h       int
	beastiary  *Beastiary
	searchTerm string
}

func NewAddCreatureWidget(name string, w, h int) *AddCreatureWidget {
	return &AddCreatureWidget{name: name, w: w, h: h, beastiary: beastiary, searchTerm: ""}
}

type AddCreatureEditor struct {
	widget *AddCreatureWidget
}

func (e *AddCreatureEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		e.widget.searchTerm += string(ch)
	} else if key == gocui.KeyBackspace || key == gocui.KeyBackspace2 {
		e.widget.searchTerm = e.widget.searchTerm[:len(e.widget.searchTerm)-1]
	}
	v.SetWritePos(0, 25)
	v.SetLine(25, "")
	v.WriteString(e.widget.searchTerm)
}

func (w *AddCreatureWidget) Layout(g *gocui.Gui) error {
	cols, lines := g.Size()

	x := int(cols / 2)
	y := int(lines / 2)

	width := int(w.w / 2)
	height := int(float32(lines)*0.8) / 2

	x0 := int(x - width)
	x1 := int(x + width)
	y0 := int(y - height)
	y1 := int(y + height)

	view, err := g.SetView(w.name, x0, y0, x1, y1, 0)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		w.searchTerm = ""
	}

	var filteredCreatures []Creature
	var index int = 1

	view.Title = "Add Creature"

	for _, creature := range w.beastiary.Creatures {

		var contains bool = false

		if w.searchTerm == "" {
			contains = true
		} else if strings.Contains(w.searchTerm, strings.ToLower(creature.Name)) {
			contains = true
		}

		if contains {
			filteredCreatures = append(filteredCreatures, creature)
			view.SetWritePos(1, index)
			fmt.Fprintf(view, "\x1b[37;2m%d\x1b[0m %s", index, creature.Name)
			index += 1
		}
	}
	view.Editable = true
	view.Editor = &AddCreatureEditor{widget: w}

	return nil

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

	g.Cursor = true

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
