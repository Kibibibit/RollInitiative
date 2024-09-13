package models

type Party struct {
	Id      string   `yaml:"-"`
	Name    string   `yaml:"name"`
	Players []Player `yaml:"players"`
}
