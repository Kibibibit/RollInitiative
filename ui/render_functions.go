package ui

import (
	"fmt"
	"strings"
	"windmills/roll_initiative/models"
)

func RenderCreatureSearchRow(c *models.Creature, colors *ColorPalette, index int, width int) string {

	line := ApplyFgColor(fmt.Sprintf(" %2d ", index), colors.FgColorDim)
	lineLength := 4
	line = fmt.Sprintf(" %s%s", line, c.Name)

	lineLength += len(c.Name)

	cr := c.CR

	for lineLength < width/2 {
		line += " "
		lineLength += 1
	}

	line += c.Type
	lineLength += len(c.Type)

	sCount := width - lineLength - 4

	oldLength := len(line)

	line += strings.Repeat(" ", sCount)

	lineLength += len(line) - oldLength

	line += cr

	lineLength += len(cr)

	for lineLength < width {
		line += " "
		lineLength += 1
	}

	return line
}
