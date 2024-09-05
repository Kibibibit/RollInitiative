package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"
)

type StatBlock struct {
	STR int `xml:"str"`
	DEX int `xml:"dex"`
	CON int `xml:"con"`
	INT int `xml:"int"`
	WIS int `xml:"wis"`
	CHA int `xml:"cha"`
}

type BeastiaryImport struct {
	Creatures []Creature `xml:"creature"`
}

type SpellListImport struct {
	Spells []Spell `xml:"spell"`
}

type SpellList = map[string]Spell
type Beastiary = map[string]Creature

type Creature struct {
	Id                    string          `xml:"id"`
	Name                  string          `xml:"name"`
	StatBlock             StatBlock       `xml:"statBlock"`
	AvgHP                 int             `xml:"avgHp"`
	AC                    int             `xml:"ac"`
	HitDice               int             `xml:"hitDice"`
	HitDiceType           int             `xml:"hitDiceType"`
	Speed                 string          `xml:"speed"`
	Type                  string          `xml:"type"`
	Size                  string          `xml:"size"`
	Alignment             string          `xml:"alignment"`
	Senses                string          `xml:"senses"`
	Languages             string          `xml:"languages"`
	CR                    string          `xml:"cr"`
	Source                string          `xml:"source"`
	Saves                 []string        `xml:"save"`
	Skills                []string        `xml:"skill"`
	DamageVulnerabilities string          `xml:"damageVulnerabilities"`
	DamageResistances     string          `xml:"damageResistances"`
	DamageImmunities      string          `xml:"damageImmunities"`
	ConditionImmunities   string          `xml:"conditionImmunities"`
	Actions               []CreatureTrait `xml:"action"`
	BonusActions          []CreatureTrait `xml:"bonusAction"`
	Reactions             []CreatureTrait `xml:"reaction"`
	LairActions           []CreatureTrait `xml:"lairAction"`
	LegendaryDescription  string          `xml:"legendaryDescription"`
	LegendaryActions      []CreatureTrait `xml:"legendaryAction"`
	SpellBook             SpellBook       `xml:"spellBook"`

	//Do spells with a list of each spell level maybe?
}

type SpellBook struct {
	Cantrips        []string `xml:"cantrip"`
	Level1          []string `xml:"level1"`
	Level2          []string `xml:"level2"`
	Level3          []string `xml:"level3"`
	Level4          []string `xml:"level4"`
	Level5          []string `xml:"level5"`
	Level6          []string `xml:"level6"`
	Level7          []string `xml:"level7"`
	Level8          []string `xml:"level8"`
	Level9          []string `xml:"level9"`
	SpellNotes      string   `xml:"spellNotes"`
	PreCombatSpells []string `xml:"precombatSpell"`
}

type Spell struct {
	Id           string          `xml:"id"`
	Name         string          `xml:"name"`
	Level        int             `xml:"level"`
	CastingTime  string          `xml:"castingTime"`
	Range        string          `xml:"range"`
	School       string          `xml:"school"`
	Duration     string          `xml:"duration"`
	Description  string          `xml:"description"`
	Ritual       bool            `xml:"ritual"`
	HigherLevels string          `xml:"higherLevels"`
	Components   SpellComponents `xml:"components"`
	Classes      []string        `xml:"class"`
}

type SpellComponents struct {
	HasVerbal   bool   `xml:"hasVerbal"`
	HasSomatic  bool   `xml:"hasSomatic"`
	HasMaterial bool   `xml:"hasMaterial"`
	Materials   string `xml:"materials"`
}

type CreatureTrait struct {
	Name        string `xml:"name"`
	Description string `xml:"desc"`
}
type IniativeEntry struct {
	creature     *Creature
	iniativeRoll int
	currentHp    int
	tag          string
}

type IniativeTracker struct {
	combatants []IniativeEntry
}

func LoadBeastiary(path string) (*BeastiaryImport, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to read creature!")

		return nil, err
	}
	defer xmlFile.Close()

	data, err := io.ReadAll(xmlFile)
	if err != nil {
		log.Fatal("Failed to parse xml for creature!")
		return nil, err
	}

	var creatures BeastiaryImport
	xml.Unmarshal(data, &creatures)

	return &creatures, nil

}

func (b *BeastiaryImport) GetCreatureByName(name string) *Creature {
	for _, c := range b.Creatures {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
