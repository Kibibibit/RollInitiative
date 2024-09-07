package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type ACJSON struct {
	Name string `json:"name"`
	AC   string `json:"Armor Class"`
}

func ReadJSON() {

	file, err := os.Open("./backup_data/new_data.json")

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	jsonList := []JSONCreature{}

	json.Unmarshal(data, &jsonList)

	yamlFile, err := os.OpenFile("./data/creatures/srd_creatures.yaml", os.O_CREATE|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer yamlFile.Close()

	yamlDict := make(map[string]Creature)

	file2, err := os.Open("./backup_data/all_srd_creatures.json")

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	acData, err := io.ReadAll(file2)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	acList := []ACJSON{}

	err = json.Unmarshal(acData, &acList)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	for _, creature := range jsonList {
		yamlCreature := JSONCreatureToCreature(creature)

		for _, acc := range acList {
			otherId := MakeId("CREATURE", acc.Name)

			if otherId == yamlCreature.Id {
				yamlCreature.AC = strings.TrimSpace(acc.AC)
			}

		}

		yamlDict[yamlCreature.Id] = yamlCreature

	}

	outData, err := yaml.Marshal(yamlDict)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	yamlFile.Write(outData)

}

type JSONCreature struct {
	Id                    string           `json:"-"`
	Name                  string           `json:"name"`
	Source                string           `json:"source,omitempty"`
	Size                  string           `json:"size"`
	Type                  string           `json:"type"`
	SubType               string           `json:"subtype"`
	Alignment             string           `json:"alignment"`
	AC                    int              `json:"ac"`
	HitDice               string           `json:"hit_dice"`
	Speed                 string           `json:"speed"`
	Stats                 []int            `json:"stats"`
	Saves                 []map[string]int `json:"saves,omitempty"`
	Skills                []map[string]int `json:"skillsaves,omitempty"`
	DamageVulnerabilities string           `json:"damage_vulnerabilities,omitempty"`
	DamageResistances     string           `json:"damage_resistances,omitempty"`
	DamageImmunities      string           `json:"damage_immunities,omitempty"`
	ConditionImmunities   string           `json:"condition_immunities,omitempty"`
	Senses                string           `json:"senses"`
	Languages             string           `json:"languages"`
	CR                    string           `json:"cr"`

	Actions              []JSONCreatureTrait `json:"actions,omitempty"`
	BonusActions         []JSONCreatureTrait `json:"bonus_actions,omitempty"`
	Reactions            []JSONCreatureTrait `json:"reactions,omitempty"`
	LairActions          []JSONCreatureTrait `json:"lair_actions,omitempty"`
	Traits               []JSONCreatureTrait `json:"traits,omitempty"`
	LegendaryDescription string              `json:"legendary_description,omitempty"`
	LegendaryActions     []JSONCreatureTrait `json:"legendary_actions,omitempty"`
	Spells               map[string]string   `json:"spells,omitempty"`
	SpellNotes           string              `json:"spellNotes"`
	PrecombatSpells      []string            `json:"precombatSpells,omitempty"`
	Slots                map[string]int      `json:"slots"`
}

type JSONCreatureTrait struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}

func convertJSONTraits(jts []JSONCreatureTrait) []CreatureTrait {
	out := []CreatureTrait{}

	for _, jt := range jts {
		out = append(out, CreatureTrait{
			Name:        jt.Name,
			Description: jt.Description,
		})
	}

	return out
}

func convertJSONSaves(jss []map[string]int) map[string]int {
	out := make(map[string]int)

	for _, entry := range jss {
		for key, value := range entry {
			out[key] = value
		}
	}

	return out
}

func JSONCreatureToCreature(jc JSONCreature) Creature {
	out := Creature{
		Id:                    MakeId("CREATURE", jc.Name),
		Name:                  jc.Name,
		Source:                jc.Source,
		Size:                  jc.Size,
		Alignment:             jc.Alignment,
		Speed:                 jc.Speed,
		Stats:                 jc.Stats,
		DamageVulnerabilities: jc.DamageVulnerabilities,
		DamageResistances:     jc.DamageResistances,
		DamageImmunities:      jc.DamageImmunities,
		ConditionImmunities:   jc.ConditionImmunities,
		Senses:                jc.Senses,
		Languages:             jc.Languages,
		CR:                    jc.CR,
		SpellNotes:            jc.SpellNotes,
	}

	out.Type = jc.Type

	if jc.SubType != "" {
		out.Type = out.Type + " (" + jc.SubType + ")"
	}

	//HitDice

	hitDiceSplit := strings.Split(jc.HitDice, "+")

	hitDiceSplit2 := strings.Split(hitDiceSplit[0], "d")

	hitDiceCount, err := strconv.ParseInt(strings.TrimSpace(hitDiceSplit2[0]), 10, 64)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	hitDiceType, err := strconv.ParseInt(strings.TrimSpace(hitDiceSplit2[1]), 10, 64)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	out.HitDice = int(hitDiceCount)
	out.HitDiceType = int(hitDiceType)

	//Saves
	out.Saves = convertJSONSaves(jc.Saves)

	//Skills
	out.Skills = convertJSONSaves(jc.Skills)

	//Actions, Bonus Actions, Reactions, Lair Actions, Legendary Actions, Traits
	out.Actions = convertJSONTraits(jc.Actions)
	out.BonusActions = convertJSONTraits(jc.BonusActions)
	out.Reactions = convertJSONTraits(jc.Reactions)
	out.LairActions = convertJSONTraits(jc.LairActions)
	out.LegendaryActions = convertJSONTraits(jc.LegendaryActions)
	out.Traits = convertJSONTraits(jc.Traits)

	//Spells
	out.Spells = make(map[int]CreatureSpells)
	for key, spells := range jc.Spells {

		level64, err := strconv.ParseInt(key, 10, 64)

		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		level := int(level64)

		spellObject := CreatureSpells{Spells: []string{}}

		if level != 0 {
			spellObject.Slots = jc.Slots[key]
		}

		spellArray := strings.Split(spells, ",")

		for _, s := range spellArray {
			spellObject.Spells = append(spellObject.Spells, MakeId("SPELL", strings.TrimSpace(s)))
		}

		out.Spells[level] = spellObject
	}

	//Precombat spells
	out.PrecombatSpells = []string{}
	for _, spellName := range jc.PrecombatSpells {
		out.PrecombatSpells = append(out.PrecombatSpells, MakeId("SPELL", spellName))
	}

	return out
}
