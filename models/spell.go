package models

import (
	"fmt"
	"strings"
)

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

func (c *Spell) RenderSearchRow(index int, width int) string {
	line := fmt.Sprintf(" %2d %s", index, c.Name)
	level := fmt.Sprintf("%d", c.Level)

	for len(line) < width/2 {
		line += " "
	}

	line += c.School

	sCount := width - len(line) - 4

	line += strings.Repeat(" ", sCount)
	line += level

	for len(line) < width {
		line += " "
	}

	return line
}
