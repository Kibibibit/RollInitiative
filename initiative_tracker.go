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
	Name        string    `xml:"name"`
	StatBlock   StatBlock `xml:"statBlock"`
	MaxHP       int       `xml:"hp"`
	AC          int       `xml:"ac"`
	HitDice     int       `xml:"hitDice"`
	HitDiceType int       `xml:"hitDiceType"`
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
