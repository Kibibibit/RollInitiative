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

type Beastiary struct {
	Creatures []Creature `xml:"creature"`
}

type Creature struct {
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
	SpellNotes            string          `xml:"spellNotes"`
	//Do spells with a list of each spell level maybe?
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

func LoadBeastiary(path string) (*Beastiary, error) {
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

	var creatures Beastiary
	xml.Unmarshal(data, &creatures)

	return &creatures, nil

}

func (b *Beastiary) GetCreatureByName(name string) *Creature {
	for _, c := range b.Creatures {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
