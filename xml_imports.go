package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func findXMLFilesInFolder(path string) ([]string, error) {
	files, err := os.ReadDir(path)

	log.Printf("Finding all xml files in %s\n", path)
	if err != nil {
		log.Println("Failed to open folder!")
		log.Fatalln(err)
		return nil, err
	}

	out := []string{}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.HasSuffix(fileName, ".xml") {
				if !strings.HasSuffix(path, "/") {
					fileName = fmt.Sprintf("/%s", fileName)
				}
				out = append(out, fmt.Sprintf("%s%s", path, fileName))
			}

		}
	}

	return out, nil
}

func readXML(path string) ([]byte, error) {
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
	return data, nil
}

func ImportSpells(path string, dict SpellDict) (SpellDict, error) {

	files, err := findXMLFilesInFolder(path)

	if err != nil {
		log.Fatalln("Failed to load spell folder!")
		log.Fatalln(err)
		return dict, err
	}

	for _, file := range files {
		dict, err = importSpellFile(file, dict)
		if err != nil {
			log.Panicf("Couldn't load file %s!\n", file)
			log.Panicln(err)
			return dict, err
		}
	}

	return dict, nil
}

func importSpellFile(path string, dict SpellDict) (SpellDict, error) {

	data, err := readXML(path)

	if err != nil {
		log.Fatalln("Failed to load spell xml file!")
		log.Fatalln(err)
		return dict, err
	}

	var spellImportList XMLSpellImportList
	err = xml.Unmarshal(data, &spellImportList)
	if err != nil {
		log.Println("Failed to unmarshal data!")
		log.Fatalln(err)
		return dict, err
	}

	log.Println("Successfully loaded spells")

	for _, item := range spellImportList.Items {
		dict[item.Id] = item
	}

	return dict, nil

}

func ImportCreatures(path string, dict CreatureDict) (CreatureDict, error) {

	files, err := findXMLFilesInFolder(path)

	if err != nil {
		log.Fatalln("Failed to load creature folder!")
		log.Fatalln(err)
		return dict, err
	}

	for _, file := range files {
		dict, err = importCreatureFile(file, dict)
		if err != nil {
			log.Panicf("Couldn't load file %s!\n", file)
			log.Panicln(err)
			return dict, err
		}
	}

	return dict, nil
}

func importCreatureFile(path string, dict CreatureDict) (CreatureDict, error) {
	data, err := readXML(path)

	if err != nil {
		log.Fatalln("Failed to load creature xml file!")
		log.Fatalln(err)
		return dict, err
	}

	var creatureImportList XMLCreatureImportList
	err = xml.Unmarshal(data, &creatureImportList)
	if err != nil {
		log.Println("Failed to unmarshal data!")
		log.Fatalln(err)
		return dict, err
	}

	log.Println("Successfully loaded creatures")

	for _, item := range creatureImportList.Items {
		dict[item.Id] = item
	}

	return dict, nil

}

func MakeId(prefix string, data string) string {
	out := strings.ToUpper(fmt.Sprintf("%s:%s", prefix, data))
	out = strings.ReplaceAll(out, " ", "_")
	out = strings.ReplaceAll(out, "(", "")
	out = strings.ReplaceAll(out, ")", "")
	out = strings.ReplaceAll(out, "-", "_")
	out = strings.ReplaceAll(out, "/", "_")
	out = strings.ReplaceAll(out, "'", "")
	out = strings.ReplaceAll(out, "’", "")
	out = strings.ReplaceAll(out, "*", "")
	return out
}
