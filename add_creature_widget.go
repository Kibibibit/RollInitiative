package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/awesome-gocui/gocui"
)

const ADD_CREATURE_NAME = "add_creature"

type AddCreatureWidget struct {
	name string
	w, h int
}

func NewAddCreatureWidget(name string, w, h int) *AddCreatureWidget {
	return &AddCreatureWidget{name: name, w: w, h: h}
}

func AddAddCreatureWidgetKeybinds(g *gocui.Gui) {

	for letter := 'A'; letter <= 'z'; letter++ {
		if err := g.SetKeybinding(ADD_CREATURE_NAME, letter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				onAddCreatureType(g, v, letter)
				return nil
			}); err != nil {
			log.Panicln(err)
		}

	}

	if err := g.SetKeybinding(ADD_CREATURE_NAME, gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			onAddCreatureType(g, v, ' ')
			return nil
		}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(ADD_CREATURE_NAME, gocui.KeyBackspace|gocui.KeyBackspace2, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			onAddCreatureBackspace(g, v)
			return nil
		}); err != nil {
		log.Panicln(err)
	}

	for number := '0'; number <= '9'; number++ {
		if err := g.SetKeybinding(ADD_CREATURE_NAME, number, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				onAddCreatureSelect(g, v, number)
				return nil
			}); err != nil {
			log.Panicln(err)
		}

	}

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

func onAddCreatureType(g *gocui.Gui, v *gocui.View, letter rune) {

	line, _ := v.Line(0)
	v.Rewind()

	searchTerm := strings.ReplaceAll(line, "🔎︎", "") + string(letter)

	AddCreatureSearch(v, searchTerm)
}

func onAddCreatureBackspace(g *gocui.Gui, v *gocui.View) {
	line, _ := v.Line(0)
	v.Rewind()

	searchTerm := strings.ReplaceAll(line, "🔎︎", "")
	if len(searchTerm) > 0 {
		searchTerm = searchTerm[0 : len(searchTerm)-1]
	}
	AddCreatureSearch(v, searchTerm)
}

func onAddCreatureSelect(g *gocui.Gui, v *gocui.View, number rune) error {

	lineNumber64, err := strconv.ParseInt(string(number), 10, 64)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	lineNumber := int(lineNumber64)

	if lineNumber == 0 {
		lineNumber = 10
	}

	lineNumber += 3

	line, err := v.Line(lineNumber)

	return switchToViewCreatureWindow(g, v, NewViewCreatureWidget(VIEW_CREATURE_NAME, 120, 50, Creature{Name: line}, ADD_CREATURE_NAME))

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
		view.Title = "Add Creature"

	}

	AddCreatureSearch(view, "")

	return nil

}

func AddCreatureSearch(v *gocui.View, searchTerm string) {
	filteredCreatures := []string{}
	var index int = 0

	v.Clear()

	lowercaseSearchTerm := strings.ToLower(searchTerm)

	w, _ := v.Size()

	v.SetWritePos(3, 3)
	fmt.Fprint(v, "\x1b[37;2mName\x1b[0m")
	v.SetWritePos(w-4, 3)
	fmt.Fprint(v, "\x1b[37;2mCR\x1b[0m")

	for _, creatureId := range creatureIds {

		creature := creatureDict[creatureId]

		var contains bool = false

		if searchTerm == "" {
			contains = true
		} else if strings.Contains(strings.ToLower(creature.Name), lowercaseSearchTerm) {
			contains = true
		}

		if contains {
			filteredCreatures = append(filteredCreatures, creatureId)
			index += 1
			if index >= 10 {
				break
			}
		}
	}

	slices.Sort(filteredCreatures)

	for index, creatureId := range filteredCreatures {
		creature := creatureDict[creatureId]
		v.SetWritePos(1, index+4)
		drawIndex := index + 1
		if index+1 == 10 {
			drawIndex = 0
		}
		fmt.Fprintf(v, "\x1b[37;2m%d\x1b[0m %s", drawIndex, creature.Name)
		v.SetWritePos(w-4, index+4)
		fmt.Fprint(v, creature.CR)
		index += 1
		if index >= 10 {
			break
		}
	}

	v.SetWritePos(0, 0)
	v.SetLine(0, "")

	fmt.Fprintf(v, "🔎︎%s", searchTerm)
	v.SetLine(1, "")
	v.SetWritePos(0, 1)
	fmt.Fprintf(v, "\x1b[37;2mFound %d results\x1b[0m", len(filteredCreatures))
}

func switchToAddCreatureWindow(g *gocui.Gui, v *gocui.View, widget *AddCreatureWidget) error {

	widget.Layout(g)

	_, err := g.SetCurrentView(ADD_CREATURE_NAME)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	return nil
}
