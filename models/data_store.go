package models

import (
	"log"
	"os"
	"strings"
	"windmills/roll_initiative/utils"
)

type DataStore struct {
	CreatureIds   []string
	CreatureNames map[string]string
	creatures     map[string]Creature
	SpellIds      []string
	SpellNames    map[string]string
	spells        map[string]Spell

	parties    map[string]Party
	PartyNames map[string]string
	PartyIds   []string

	players     map[string]Player
	PlayerNames map[string]string
	PlayerIds   []string

	IniativeEntries map[int]*IniativeEntry
}

func MakeDataStore(dataFolders []string) *DataStore {
	out := DataStore{}

	spellDict := make(map[string]Spell)
	creatureDict := make(map[string]Creature)
	partyDict := make(map[string]Party)
	playerDict := make(map[string]Player)

	var err error

	for _, folder := range dataFolders {
		folder = strings.TrimSpace(folder)
		if !strings.HasSuffix(folder, "/") {
			folder = folder + "/"
		}

		spellFolder := folder + "spells"
		creatureFolder := folder + "creatures"
		partyFolder := folder + "parties"

		spellDict, err = ImportSpells(spellFolder, spellDict)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		partyDict, err = ImportParties(partyFolder, partyDict)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		creatureDict, err = ImportCreatures(creatureFolder, creatureDict)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}

	spellIds := make([]string, len(spellDict))
	spellNames := make(map[string]string)
	var i int = 0
	for k := range spellDict {
		s := spellDict[k]
		s.Id = k
		spellDict[k] = s
		spellIds[i] = k
		spellNames[k] = s.Name
		i++
	}

	creatureIds := make([]string, len(creatureDict))
	creatureNames := make(map[string]string)

	i = 0
	for k := range creatureDict {
		c := creatureDict[k]
		c.Id = k
		creatureDict[k] = c
		creatureIds[i] = k
		creatureNames[k] = c.Name

		i++
	}

	partyIds := make([]string, len(partyDict))
	partyNames := make(map[string]string)
	playerIds := []string{}
	playerNames := make(map[string]string)

	i = 0
	for partyKey := range partyDict {
		party := partyDict[partyKey]
		partyIds[i] = partyKey
		partyNames[partyKey] = party.Name

		for _, player := range party.Players {
			playerDict[player.Id] = player
			playerIds = append(playerIds, player.Id)
			playerNames[player.Id] = player.Name
		}

		i++
	}

	out.SpellIds = spellIds
	out.CreatureIds = creatureIds
	out.creatures = creatureDict
	out.spells = spellDict
	out.CreatureNames = creatureNames
	out.SpellNames = spellNames

	out.parties = partyDict
	out.PartyIds = partyIds
	out.PartyNames = partyNames
	out.players = playerDict
	out.PlayerIds = playerIds
	out.PlayerNames = playerNames

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

func (d *DataStore) GetParty(id string) *Party {
	if party, ok := d.parties[id]; ok {
		return &party
	}
	return nil
}

func (d *DataStore) GetPlayer(id string) *Player {
	if player, ok := d.players[id]; ok {
		return &player
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
		MaxHp:        hp,
		Tag:          strings.TrimSpace(tag),
		IniativeRoll: intRoll,
		DexScore:     creature.GetDex(),
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

func (d *DataStore) NewPartyEntries(partyId string) []*IniativeEntry {

	party := d.GetParty(partyId)

	if party == nil {
		log.Fatalln("Failed to find party!")
	}

	for _, player := range party.Players {

		entry := &IniativeEntry{
			CreatureId: player.Id,
			IsPlayer:   true,
			DexScore:   player.DexScore,
		}

		entryId := 0

		for _, existingEntry := range d.IniativeEntries {
			if entryId <= existingEntry.EntryId {
				entryId = existingEntry.EntryId + 1
			}

			entry.EntryId = entryId
		}

		d.IniativeEntries[entryId] = entry

		if len(player.Familiars) > 0 {
			for _, familiar := range player.Familiars {
				d.NewCreatureEntry(familiar.CreatureId, familiar.Name, false)
			}
		}

	}

	return nil
}

func (d *DataStore) DeleteCreatureEntry(entryId int) {
	delete(d.IniativeEntries, entryId)
}
