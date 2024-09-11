package ui

import (
	"fmt"
	"strings"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"
)

func RenderCreatureSearchRow(c *models.Creature, colors *ColorPalette, index int, width int) string {

	line := ApplyFgColor(fmt.Sprintf(" %2d ", index), colors.FgColorDim)

	line = fmt.Sprintf(" %s%s", line, c.Name)

	cr := c.CR

	for utils.StringDrawLength(line) < width/2 {
		line += " "
	}

	line += c.Type

	sCount := width - utils.StringDrawLength(line) - 4

	line += strings.Repeat(" ", sCount)

	line += cr

	for utils.StringDrawLength(line) < width {
		line += " "
	}

	return line
}

func RenderSpellSearchRow(s *models.Spell, colors *ColorPalette, index int, width int) string {
	line := ApplyFgColor(fmt.Sprintf(" %2d ", index), colors.FgColorDim)
	line = fmt.Sprintf(" %s%s", line, s.Name)

	level := fmt.Sprintf("%d", s.Level)

	for utils.StringDrawLength(line) < width/2 {
		line += " "
	}

	line += s.School

	sCount := width - utils.StringDrawLength(line) - 4

	line += strings.Repeat(" ", sCount)
	line += level

	for utils.StringDrawLength(line) < width {
		line += " "
	}

	return line
}
