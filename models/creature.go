package models

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
