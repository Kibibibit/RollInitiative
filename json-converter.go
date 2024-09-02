package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func ImportMonsters() error {
	jsonFile, err := os.Open("./data/data.json")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to read creature!")

		return err
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Failed to parse xml for creature!")
		return err
	}

	err = os.Remove("./data/srd_creatures.xml")
	if err != nil {
		log.Fatal(err)

		return err
	}
	xmlFile, err := os.Create("./data/srd_creatures.xml")

	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to load creature!")

		return err
	}

	defer xmlFile.Close()
	var creatures []MonsterImport
	json.Unmarshal(data, &creatures)

	var beastiary Beastiary
	beastiary.Creatures = []Creature{}

	for _, imp := range creatures {

		creature := MonsterImportToCreature(&imp)

		beastiary.Creatures = append(beastiary.Creatures, *creature)

	}

	b, err := xml.MarshalIndent(beastiary, "", "	")

	dataString := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + string(b)
	dataString = strings.ReplaceAll(dataString, "Beastiary", "creatures")
	dataString = strings.ReplaceAll(dataString, "&#39;", "'")

	if err != nil {
		log.Fatal("Failed to convert creature")
		return err
	}
	_, err = xmlFile.WriteString(dataString)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func MonsterImportToCreature(in *MonsterImport) *Creature {
	var out Creature = Creature{}

	out.Name = in.Name
	out.AvgHP = in.HP
	out.AC = in.AC
	out.Speed = in.Speed
	out.Alignment = in.Alignment
	out.Size = in.Size
	out.Type = in.Type
	if len(in.Subtype) > 0 {
		out.Type = fmt.Sprintf("%s (%s)", out.Type, in.Subtype)
	}

	out.Senses = in.Senses
	out.Languages = in.Languages
	out.CR = in.CR

	out.StatBlock.STR = in.Stats[0]
	out.StatBlock.DEX = in.Stats[1]
	out.StatBlock.CON = in.Stats[2]
	out.StatBlock.INT = in.Stats[3]
	out.StatBlock.WIS = in.Stats[4]
	out.StatBlock.CHA = in.Stats[5]

	out.Source = in.Source

	hitDiceString := strings.Split(in.HitDice, " ")[0]
	hitDiceValues := strings.Split(hitDiceString, "d")

	hitDiceCount, err := strconv.ParseInt(hitDiceValues[0], 10, 64)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	out.HitDice = int(hitDiceCount)

	hitDiceType, err := strconv.ParseInt(hitDiceValues[1], 10, 64)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	out.HitDiceType = int(hitDiceType)

	out.DamageImmunities = in.DamageImmunities
	out.DamageResistances = in.DamageResistances
	out.DamageVulnerabilities = in.DamageVulnerabilities
	out.ConditionImmunities = in.ConditionImmunities

	out.Saves = convertWeirdIntStringMap(in.Saves)
	out.Skills = convertWeirdIntStringMap(in.SkillSaves)

	out.Actions = convertTraits(in.Actions)
	out.BonusActions = convertTraits(in.BonusActions)
	out.Reactions = convertTraits(in.Reactions)
	out.LairActions = convertTraits(in.LairActions)
	out.LegendaryActions = convertTraits(in.LegendaryActions)
	out.LegendaryDescription = in.LegendaryDescription

	out.SpellNotes = in.SpellNotes

	return &out
}

func convertWeirdIntStringMap(v []map[string]int) []string {
	out := []string{}

	for _, entry := range v {

		for name, value := range entry {

			valueString := fmt.Sprintf("+%d", value)
			if value < 0 {
				valueString = fmt.Sprintf("%d", value)
			}
			if value == 0 {
				continue
			}

			out = append(out, fmt.Sprintf("%s: %s", name, valueString))
		}
	}

	return out
}

func convertTraits(v []MonsterImportTrait) []CreatureTrait {
	var out []CreatureTrait

	for _, i := range v {
		out = append(out, CreatureTrait{
			Name:        i.Name,
			Description: i.Description,
		})
	}

	return out
}

type MonsterImportTrait struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type MonsterImportSpell = []map[string]string

type MonsterImport struct {
	Name                  string               `json:"name"`
	Source                string               `json:"source"`
	Size                  string               `json:"size"`
	Type                  string               `json:"type"`
	Subtype               string               `json:"subtype"`
	Alignment             string               `json:"alignment"`
	AC                    int                  `json:"ac"`
	HP                    int                  `json:"hp"`
	HitDice               string               `json:"hit_dice"`
	Speed                 string               `json:"speed"`
	Stats                 []int                `json:"stats"`
	Saves                 []map[string]int     `json:"saves"`
	SkillSaves            []map[string]int     `json:"skillsaves"`
	DamageVulnerabilities string               `json:"damage_vulnerabilities"`
	DamageResistances     string               `json:"damage_resistances"`
	DamageImmunities      string               `json:"damage_immunities"`
	ConditionImmunities   string               `json:"condition_immunities"`
	Senses                string               `json:"senses"`
	Languages             string               `json:"languages"`
	CR                    string               `json:"cr"`
	Traits                []MonsterImportTrait `json:"traits"`
	Actions               []MonsterImportTrait `json:"actions"`
	BonusActions          []MonsterImportTrait `json:"bonus_actions"`
	Reactions             []MonsterImportTrait `json:"reactions"`
	LairActions           []MonsterImportTrait `json:"lair_actions"`
	LegendaryActions      []MonsterImportTrait `json:"legendary_actions"`
	LegendaryDescription  string               `json:"legendary_description"`
	SpellNotes            string               `json:"spellsNotes"`
	Spells                MonsterImportSpell   `json:"spells"`
}
