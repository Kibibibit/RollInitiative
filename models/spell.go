package models

type Spell struct {
	Id           string   `yaml:"-"`
	Name         string   `yaml:"name"`
	Level        int      `yaml:"level"`
	CastingTime  string   `yaml:"castingTime"`
	Range        string   `yaml:"range"`
	School       string   `yaml:"school"`
	Duration     string   `yaml:"duration"`
	Description  string   `yaml:"description"`
	Ritual       bool     `yaml:"ritual"`
	HigherLevels string   `yaml:"higherLevels,omitempty"`
	Components   string   `yaml:"components"`
	Materials    string   `yaml:"materials,omitempty"`
	Classes      []string `yaml:"class,omitempty"`
	Source       string   `yaml:"source"`
}

func (s Spell) GetId() string {
	return s.Id
}

func (s Spell) GetName() string {
	return s.Name
}
