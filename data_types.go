package main

type SpellDict = map[string]Spell
type CreatureDict = map[string]Creature

type Creature struct {
	Id                    string          `xml:"id"`
	Name                  string          `xml:"name"`
	StatBlock             StatBlock       `xml:"statBlock"`
	AvgHP                 int             `xml:"avgHp"`
	AC                    string          `xml:"ac"`
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
	Traits                []CreatureTrait `xml:"trait"`
	LegendaryDescription  string          `xml:"legendaryDescription"`
	LegendaryActions      []CreatureTrait `xml:"legendaryAction"`
	SpellBook             SpellBook       `xml:"spellBook"`

	//Do spells with a list of each spell level maybe?
}

type StatBlock struct {
	STR int `xml:"str"`
	DEX int `xml:"dex"`
	CON int `xml:"con"`
	INT int `xml:"int"`
	WIS int `xml:"wis"`
	CHA int `xml:"cha"`
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
