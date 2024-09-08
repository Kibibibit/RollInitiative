package models

import (
	"log"
	"os"
)

type DataStore struct {
	CreatureIds   []string
	CreatureNames map[string]string
	creatures     map[string]Creature
	SpellIds      []string
	spells        map[string]Spell
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
	creatureNames := make(map[string]string)

	i = 0
	for k := range creatureDict {
		creatureIds[i] = k
		creatureNames[k] = creatureDict[k].Name
		i++
	}

	out.SpellIds = spellIds
	out.CreatureIds = creatureIds
	out.creatures = creatureDict
	out.spells = spellDict
	out.CreatureNames = creatureNames

	return &out

}

func (d *DataStore) GetSpell(id string) *Spell {
	if spell, ok := d.spells[id]; ok {
		return &spell
	}
	return nil
}

func (d *DataStore) GetCreature(id string) *Creature {
	if creature, ok := d.creatures[id]; ok {
		return &creature
	}
	return nil
}
