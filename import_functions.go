package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func findYAMLFilesInFolder(path string) ([]string, error) {
	files, err := os.ReadDir(path)

	log.Printf("Finding all yaml files in %s\n", path)
	if err != nil {
		log.Println("Failed to open folder!")
		log.Fatalln(err)
		return nil, err
	}

	out := []string{}

	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(path, "/") {
			fileName = fmt.Sprintf("/%s", fileName)
		}
		if !file.IsDir() {

			if strings.HasSuffix(fileName, ".yaml") {

				out = append(out, fmt.Sprintf("%s%s", path, fileName))
			}

		} else {
			subfile, err := findYAMLFilesInFolder(fmt.Sprintf("%s%s", path, fileName))

			if err != nil {
				log.Println("Failed to open folder!")
				log.Fatalln(err)
				return nil, err
			}
			out = append(out, subfile...)
		}
	}

	return out, nil
}

func readYAML(path string) ([]byte, error) {
	yamlFile, err := os.Open(path)

	log.Printf("Trying to open file %s\n", path)

	if err != nil {
		log.Println("Failed to open file!")
		log.Fatalln(err)
		return nil, err
	}

	defer yamlFile.Close()

	log.Printf("Trying to read data from file %s\n", path)

	data, err := io.ReadAll(yamlFile)

	if err != nil {
		log.Println("Failed to read data!")
		log.Fatalln(err)
		return nil, err
	}
	return data, nil
}

func unmarshalYamlData[V any](path string, dict map[string]V) (map[string]V, error) {

	data, err := readYAML(path)

	if err != nil {
		log.Fatalln("Failed to load xml file!")
		log.Fatalln(err)
		return dict, err
	}

	var newItems map[string]V
	err = yaml.Unmarshal(data, &newItems)
	if err != nil {
		log.Println("Failed to unmarshal data!")
		log.Fatalln(err)
		return dict, err
	}

	log.Printf("Loaded %d new items from %s\n", len(newItems), path)

	for key, value := range newItems {
		dict[key] = value
	}

	return dict, nil
}

func importGenericData[V any](path string, dict map[string]V) (map[string]V, error) {
	files, err := findYAMLFilesInFolder(path)

	if err != nil {
		log.Fatalf("Failed to load folder %s!\n", path)
		log.Fatalln(err)
		return dict, err
	}

	for _, file := range files {
		dict, err = unmarshalYamlData(file, dict)
		if err != nil {
			log.Panicf("Couldn't load file %s!\n", file)
			log.Panicln(err)
			return dict, err
		}
	}

	return dict, nil
}

func ImportSpells(path string, dict SpellDict) (SpellDict, error) {

	dict, err := importGenericData(path, dict)

	if err != nil {
		log.Fatalln("Failed to import spells!")
		log.Fatalln(err)
		return nil, err
	}

	return dict, nil
}

func ImportCreatures(path string, dict CreatureDict) (CreatureDict, error) {

	dict, err := importGenericData(path, dict)

	if err != nil {
		log.Fatalln("Failed to import creatures!")
		log.Fatalln(err)
		return nil, err
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
