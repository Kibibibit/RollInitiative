package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"
)

func LoadBeastiary(path string) (*XMLCreatureImportList, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Failed to read creature!")

		return nil, err
	}
	defer xmlFile.Close()

	data, err := io.ReadAll(xmlFile)
	if err != nil {
		log.Fatal("Failed to parse xml for creature!")
		return nil, err
	}

	var creatures XMLCreatureImportList
	xml.Unmarshal(data, &creatures)

	return &creatures, nil

}

func (b *XMLCreatureImportList) GetCreatureByName(name string) *Creature {
	for _, c := range b.Items {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
