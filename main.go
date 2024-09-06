package main

import (
	"errors"
	"log"
	"os"
	"slices"

	"github.com/awesome-gocui/gocui"
)

const ADD_CREATURE_NAME = "add_creature"

var (
	views        = []string{}
	spellDict    SpellDict
	creatureDict CreatureDict
	spellIds     []string
	creatureIds  []string
)

func main() {

	spellDict = make(SpellDict)
	creatureDict = make(CreatureDict)

	spellDict, err := ImportSpells("./data/spells", spellDict)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	creatureDict, err = ImportCreatures("./data/creatures", creatureDict)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	spellIds = make([]string, len(spellDict))
	creatureIds = make([]string, len(creatureDict))

	var i int = 0
	for k := range spellDict {
		spellIds[i] = k
		i++
	}
	i = 0
	for k := range creatureDict {
		creatureIds[i] = k
		i++
	}

	slices.Sort(spellIds)
	slices.Sort(creatureIds)

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	table := NewMainTableWidget(MAIN_TABLE_NAME, 0, 0)

	addCreature := NewAddCreatureWidget(ADD_CREATURE_NAME, 50, 16)

	g.SetManager(table)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	AddMainTableWidgetKeybinds(g, addCreature)
	AddAddCreatureWidgetKeybinds(g)

	spellDict["1"] = Spell{}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

}
