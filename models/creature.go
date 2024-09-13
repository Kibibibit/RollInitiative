package models

import "math"

var crToPBTable = map[string]int{
	"0":   2,
	"1/8": 2,
	"1/4": 2,
	"1/2": 2,
	"1":   2,
	"2":   2,
	"3":   2,
	"4":   2,
	"5":   3,
	"6":   3,
	"7":   3,
	"8":   3,
	"9":   4,
	"10":  4,
	"11":  4,
	"12":  4,
	"13":  5,
	"14":  5,
	"15":  5,
	"16":  5,
	"17":  6,
	"18":  6,
	"19":  6,
	"20":  6,
	"21":  7,
	"22":  7,
	"23":  7,
	"24":  7,
	"25":  8,
	"26":  8,
	"27":  8,
	"28":  8,
	"29":  9,
	"30":  9,
}

var crToXPTable = map[string]int{
	"0":   10,
	"1/8": 25,
	"1/4": 50,
	"1/2": 100,
	"1":   200,
	"2":   450,
	"3":   700,
	"4":   1100,
	"5":   1800,
	"6":   2300,
	"7":   2900,
	"8":   3900,
	"9":   5000,
	"10":  5900,
	"11":  7200,
	"12":  8400,
	"13":  10000,
	"14":  11500,
	"15":  13000,
	"16":  15000,
	"17":  18000,
	"18":  20000,
	"19":  22000,
	"20":  25000,
	"21":  33000,
	"22":  41000,
	"23":  50000,
	"24":  62000,
	"25":  75000,
	"26":  90000,
	"27":  105000,
	"28":  120000,
	"29":  135000,
	"30":  155000,
}

func XPFromCR(cr string) (xp int) {
	if val, ok := crToXPTable[cr]; ok {
		return val
	}
	return 0
}

func PBFromCR(cr string) (pb int) {
	if val, ok := crToPBTable[cr]; ok {
		return val
	}
	return 0
}

type Creature struct {
	Id                    string                 `yaml:"-"`
	Name                  string                 `yaml:"name"`
	Stats                 []int                  `yaml:"stats"`
	AC                    string                 `yaml:"ac"`
	HitDice               int                    `yaml:"hitDice"`
	HitDiceType           int                    `yaml:"hitDiceType"`
	Speed                 string                 `yaml:"speed"`
	Type                  string                 `yaml:"type"`
	Size                  string                 `yaml:"size"`
	Alignment             string                 `yaml:"alignment,omitempty"`
	Senses                string                 `yaml:"senses,omitempty"`
	Languages             string                 `yaml:"languages,omitempty"`
	CR                    string                 `yaml:"cr"`
	Source                string                 `yaml:"source,omitempty"`
	Saves                 map[string]int         `yaml:"saves,omitempty"`
	Skills                map[string]int         `yaml:"skills,omitempty"`
	DamageVulnerabilities string                 `yaml:"damageVulnerabilities,omitempty"`
	DamageResistances     string                 `yaml:"damageResistances,omitempty"`
	DamageImmunities      string                 `yaml:"damageImmunities,omitempty"`
	ConditionImmunities   string                 `yaml:"conditionImmunities,omitempty"`
	Actions               []CreatureTrait        `yaml:"actions,omitempty"`
	BonusActions          []CreatureTrait        `yaml:"bonusActions,omitempty"`
	Reactions             []CreatureTrait        `yaml:"reactions,omitempty"`
	LairActions           []CreatureTrait        `yaml:"lairActions,omitempty"`
	Traits                []CreatureTrait        `yaml:"traits,omitempty"`
	LegendaryActions      []CreatureTrait        `yaml:"legendaryActions,omitempty"`
	LegendaryDescription  string                 `yaml:"legendaryDescription,omitempty"`
	SpellNotes            string                 `yaml:"spellNotes,omitempty"`
	Spells                map[int]CreatureSpells `yaml:"spells,omitempty"`
	PrecombatSpells       []string               `yaml:"precombatSpells,omitempty"`
}

type CreatureTrait struct {
	Name        string `yaml:"name"`
	Description string `yaml:"desc"`
}

type CreatureSpells struct {
	Slots  int      `yaml:"slots,omitempty"`
	Spells []string `yaml:"spells"`
}

func (c *Creature) GetMod(x int) int {
	return int(math.Floor(float64(x-10.0) / 2.0))
}

func (c *Creature) GetStr() int {
	return c.Stats[0]
}

func (c *Creature) GetDex() int {
	return c.Stats[1]
}

func (c *Creature) GetCon() int {
	return c.Stats[2]
}

func (c *Creature) GetInt() int {
	return c.Stats[3]
}

func (c *Creature) GetWis() int {
	return c.Stats[4]
}

func (c *Creature) GetCha() int {
	return c.Stats[5]
}

func (c *Creature) GetStrMod() int {
	return c.GetMod(c.Stats[0])
}

func (c *Creature) GetDexMod() int {
	return c.GetMod(c.Stats[1])
}

func (c *Creature) GetConMod() int {
	return c.GetMod(c.Stats[2])
}

func (c *Creature) GetIntMod() int {
	return c.GetMod(c.Stats[3])
}

func (c *Creature) GetWisMod() int {
	return c.GetMod(c.Stats[4])
}

func (c *Creature) GetChaMod() int {
	return c.GetMod(c.Stats[5])
}
