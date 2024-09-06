package main

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

type SpellComponents struct {
	HasVerbal   bool   `xml:"hasVerbal"`
	HasSomatic  bool   `xml:"hasSomatic"`
	HasMaterial bool   `xml:"hasMaterial"`
	Materials   string `xml:"materials"`
}
