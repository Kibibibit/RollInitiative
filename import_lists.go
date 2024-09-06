package main

type XMLCreatureImportList struct {
	Items []Creature `xml:"creature"`
}

type XMLSpellImportList struct {
	Items []Spell `xml:"spell"`
}
