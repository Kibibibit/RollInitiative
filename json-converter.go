package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func ImportMonsters(path string, spellList *SpellDict) ([]Creature, error) {
	xmlFile, err := os.Open(path)

	log.Printf("Trying to open file %s\n", path)

	if err != nil {
		log.Println("Failed to open file!")
		log.Fatalln(err)
		return nil, err
	}

	defer xmlFile.Close()

	log.Printf("Trying to read data from file %s\n", path)

	data, err := io.ReadAll(xmlFile)

	if err != nil {
		log.Println("Failed to read data!")
		log.Fatalln(err)
		return nil, err
	}

	var creatureImportList XMLCreatureImportList

	log.Printf("Attempting to unmarshal creature data from %s\n", path)

	err = xml.Unmarshal(data, &creatureImportList)
	if err != nil {
		log.Println("Failed to unmarshal data!")
		log.Fatalln(err)
		return nil, err
	}

	return creatureImportList.Items, nil

}

func ConvertSpellsToXML() error {
	jsonFile, err := os.Open("./data/spells.json")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to read spell json!")

		return err
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Failed to parse json for spell!")
		return err
	}

	err = os.Remove("./data/srd_spells.xml")
	if err != nil {

		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
			return err
		}
	}
	xmlFile, err := os.Create("./data/srd_spells.xml")

	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to load spell xml!")

		return err
	}

	defer xmlFile.Close()

	var spellImportList []SpellImport

	json.Unmarshal(data, &spellImportList)

	var spellList XMLSpellImportList

	spellList.Items = []Spell{}

	for _, imp := range spellImportList {

		spell := SpellImportToSpell(&imp)

		spellList.Items = append(spellList.Items, *spell)

	}

	s, err := xml.MarshalIndent(spellList, "", "	")

	dataString := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + string(s)
	dataString = strings.ReplaceAll(dataString, "SpellListImport", "spells")

	if err != nil {
		log.Fatal("Failed to convert spell")
		return err
	}
	_, err = xmlFile.WriteString(dataString)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func ConvertCreaturesToXML(spellList SpellDict) error {
	jsonFile, err := os.Open("./data/new_data.json")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to read creature!")

		return err
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Failed to parse json for creature!")
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

	var beastiary XMLCreatureImportList
	beastiary.Items = []Creature{}

	for _, imp := range creatures {

		creature := MonsterImportToCreature(&imp, spellList)

		beastiary.Items = append(beastiary.Items, *creature)

	}

	b, err := xml.MarshalIndent(beastiary, "", "	")

	dataString := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + string(b)
	dataString = strings.ReplaceAll(dataString, "BeastiaryImport", "creatures")

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

func SpellImportToSpell(in *SpellImport) *Spell {
	var out Spell = Spell{}

	out.Id = MakeId("SPELL", in.Name)
	out.Name = in.Name
	out.Name = strings.ReplaceAll(out.Name, "’", "'")

	out.Level = 0
	if in.Level != "cantrip" {
		level, err := strconv.ParseInt(in.Level, 10, 64)
		if err != nil {
			log.Println(in)
			log.Panic(err)
			os.Exit(1)
		}
		out.Level = int(level)
	}

	out.School = in.School
	out.Range = in.Range
	out.Description = in.Description
	out.CastingTime = in.CastingTime
	out.Duration = in.Duration
	out.Ritual = in.Ritual
	out.HigherLevels = in.HigherLevels

	out.Components = SpellComponents{}
	out.Components.HasMaterial = in.Components.HasMaterial
	out.Components.HasSomatic = in.Components.HasSomatic
	out.Components.HasVerbal = in.Components.HasVerbal
	out.Components.Materials = strings.Join(in.Components.Materials, ", ")
	out.Classes = in.Classes

	return &out
}

func MonsterImportToCreature(in *MonsterImport, spells SpellDict) *Creature {
	var out Creature = Creature{}

	out.Id = MakeId("CREATURE", in.Name)
	out.Name = in.Name
	out.AvgHP = in.HP
	out.AC = fmt.Sprintf("%d", in.AC)
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

	out.SpellBook = SpellBook{}
	out.SpellBook.SpellNotes = in.SpellNotes

	for level, spell := range in.Spells {
		spellNames := strings.Split(spell, ",")
		for _, name := range spellNames {
			name = strings.Trim(name, " ")
			spellId := MakeId("SPELL", name)

			if _, ok := spells[spellId]; ok {
				switch level {
				case "0":
					out.SpellBook.Cantrips = append(out.SpellBook.Cantrips, spellId)
					break
				case "1":
					out.SpellBook.Level1 = append(out.SpellBook.Level1, spellId)
					break
				}
			} else {
				log.Fatalf("Couldnt find spell with id %s for creature with id %s", spellId, out.Id)
			}

		}
	}

	for _, spell := range in.PreCombatSpells {
		spellId := MakeId("SPELL", spell)
		if _, ok := spells[spellId]; ok {
			out.SpellBook.PreCombatSpells = append(out.SpellBook.PreCombatSpells, spellId)
		} else {
			log.Fatalf("Couldnt find pre-combat spell with id %s for creature with id %s", spellId, out.Id)
		}
	}
	out.Traits = convertTraits(in.Traits)
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

type MonsterImportSpell = map[string]string

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
	SpellNotes            string               `json:"spellNotes"`
	Spells                MonsterImportSpell   `json:"spells"`
	PreCombatSpells       []string             `json:"precombatSpells"`
}

type SpellImport struct {
	CastingTime  string                `json:"casting_time"`
	Name         string                `json:"name"`
	Level        string                `json:"level"`
	Range        string                `json:"range"`
	School       string                `json:"school"`
	Duration     string                `json:"duration"`
	Description  string                `json:"description"`
	Ritual       bool                  `json:"ritual"`
	HigherLevels string                `json:"higher_levels"`
	Components   SpellComponentsImport `json:"components"`
	Classes      []string              `json:"classes"`
}

type SpellComponentsImport struct {
	HasVerbal   bool     `json:"verbal"`
	HasSomatic  bool     `json:"somatic"`
	HasMaterial bool     `json:"material"`
	Materials   []string `json:"materials_needed"`
}
