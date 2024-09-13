package models

type Player struct {
	Id        string     `yaml:"id"`
	Name      string     `yaml:"name"`
	DexScore  int        `yaml:"dex"`
	Familiars []Familiar `yaml:"familiars,omitempty"`
}
