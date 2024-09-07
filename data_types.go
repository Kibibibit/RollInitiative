package main

type IniativeEntry struct {
	creature     *Creature
	iniativeRoll int
	currentHp    int
	tag          string
}

type IniativeTracker struct {
	combatants []IniativeEntry
}
