package models

type IniativeEntry struct {
	CreatureId   string
	IsPlayer     bool
	Hp           int
	MaxHp        int
	IniativeRoll int
	Statuses     string
	Tag          string
	EntryId      int
	DexScore     int
}
