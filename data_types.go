package main

type SpellDict = map[string]Spell
type CreatureDict = map[string]Creature

type IniativeEntry struct {
	creature     *Creature
	iniativeRoll int
	currentHp    int
	tag          string
}

type IniativeTracker struct {
	combatants []IniativeEntry
}
