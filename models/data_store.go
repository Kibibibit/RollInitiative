package models

import (
	"log"
	"os"
	"slices"
)

type DataStore struct {
	creatureIds []string
	creatures   map[string]Creature
	spellIds    []string
	spells      map[string]Spell
}

func MakeDataStore() *DataStore {
	out := DataStore{}

	spellDict := make(map[string]Spell)
	creatureDict := make(map[string]Creature)

	spellDict, err := ImportSpells("./data/spells", spellDict)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	spellIds := make([]string, len(spellDict))
	var i int = 0
	for k := range spellDict {
		spellIds[i] = k
		i++
	}

	creatureDict, err = ImportCreatures("./data/creatures", creatureDict)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	creatureIds := make([]string, len(creatureDict))

	i = 0
	for k := range creatureDict {
		creatureIds[i] = k
		i++
	}

	slices.Sort(spellIds)
	slices.Sort(creatureIds)

	out.spellIds = spellIds
	out.creatureIds = creatureIds
	out.creatures = creatureDict
	out.spells = spellDict

	return &out

}
