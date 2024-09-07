package main

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
	LegendaryDescription  string                 `yaml:"legendaryDescription,omitempty"`
	LegendaryActions      []CreatureTrait        `yaml:"legendaryActions,omitempty"`
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
