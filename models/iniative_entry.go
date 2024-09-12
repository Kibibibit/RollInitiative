package models

type IniativeEntry struct {
	CreatureId   string
	IsPlayer     bool
	Hp           int
	IniativeRoll int
	Statuses     string
	Tag          string
	EntryId      int
}
