package models

import (
	"log"
	"os"
	"windmills/roll_initiative/utils"
)

type DataStore struct {
	CreatureIds   []string
	CreatureNames map[string]string
	creatures     map[string]Creature
	SpellIds      []string
	SpellNames    map[string]string
	spells        map[string]Spell

	IniativeEntries map[int]*IniativeEntry
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
	spellNames := make(map[string]string)
	var i int = 0
	for k := range spellDict {
		s := spellDict[k]
		spellIds[i] = k
		spellNames[k] = s.Name
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
		c := creatureDict[k]
		creatureIds[i] = k
		creatureNames[k] = c.Name

		i++
	}

	out.SpellIds = spellIds
	out.CreatureIds = creatureIds
	out.creatures = creatureDict
	out.spells = spellDict
	out.CreatureNames = creatureNames
	out.SpellNames = spellNames

	out.IniativeEntries = make(map[int]*IniativeEntry)

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

func (d *DataStore) NewCreatureEntry(cId string, tag string, rollHp bool) *IniativeEntry {

	creature := d.GetCreature(cId)
	if creature == nil {
		log.Fatalln("Failed to find creature!")
	}

	var hp int
	conMod := creature.GetConMod()
	if rollHp {
		hp = utils.RollDice(creature.HitDice, creature.HitDiceType)
	} else {
		hp = utils.AverageDiceRoll(creature.HitDice, creature.HitDiceType)
	}
	hp += conMod * creature.HitDice

	intRoll := utils.RollDice(1, 20)
	intRoll += creature.GetDexMod()

	entry := &IniativeEntry{
		CreatureId:   cId,
		IsPlayer:     false,
		Statuses:     "",
		Hp:           hp,
		Tag:          tag,
		IniativeRoll: intRoll,
	}

	entryId := 0

	for _, entry := range d.IniativeEntries {
		if entryId <= entry.EntryId {
			entryId = entry.EntryId + 1
		}
	}

	entry.EntryId = entryId

	d.IniativeEntries[entryId] = entry

	return entry

}

func (d *DataStore) DeleteCreatureEntry(entryId int) {
	delete(d.IniativeEntries, entryId)
}
