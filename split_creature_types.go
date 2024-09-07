package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func CreateCreatureFiles() {

	creatureDicts := make(map[string]SpellDict)

	for _, cid := range spellIds {
		c := spellDict[cid]

		cType := strings.ToLower(strings.Split(c.School, " ")[0])

		_, ok := creatureDicts[cType]

		typeDict := make(SpellDict)
		if ok {
			for key, value := range creatureDicts[cType] {
				typeDict[key] = value
			}

		}
		typeDict[cid] = c

		creatureDicts[cType] = typeDict

	}

	log.Println(creatureDicts)

	for key, value := range creatureDicts {
		filename := fmt.Sprintf("./data/test/srd_spells_%s.yaml", key)

		CreateCreatureFile(filename, value)
	}

}

func CreateCreatureFile(filename string, creatures SpellDict) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.Remove(filename)
	}

	log.Println(filename)

	f, err := os.Create(filename)
	if err != nil {
		log.Panicln(err)
	}

	defer f.Close()

	byteData, err := yaml.Marshal(creatures)
	if err != nil {
		log.Panicln(err)
	}

	_, err = f.Write(byteData)
	if err != nil {
		log.Panicln(err)
	}
}
