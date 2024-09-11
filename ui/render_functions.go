package ui

import (
	"fmt"
	"strings"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

func DrawText(view *gocui.View, colW int, maxH int, text string, drawX, drawY int) (int, int) {
	view.SetWritePos(drawX, drawY)

	lines := []string{}

	baseLines := strings.Split(text, "\n")

	for _, baseLine := range baseLines {

		if utils.StringDrawLength(baseLine) < colW-2 {
			lines = append(lines, baseLine)
		} else {
			words := strings.Split(baseLine, " ")

			newLine := ""
			for len(words) > 0 {
				nextWord := words[0]

				if utils.StringDrawLength(newLine)+utils.StringDrawLength(nextWord)+1 < colW-2 {
					if len(newLine) == 0 {
						newLine = nextWord
					} else {
						newLine = fmt.Sprintf("%s %s", newLine, nextWord)
					}

					words = words[1:len(words)]
				} else {
					lines = append(lines, newLine)
					newLine = ""
				}
			}
			if len(newLine) > 0 {
				lines = append(lines, newLine)
			}

		}
	}

	for _, line := range lines {

		view.SetWritePos(drawX, drawY)
		fmt.Fprint(view, line)
		drawY += 1
		if drawY > maxH-1 {
			drawY = 1
			drawX += colW
		}

	}

	return drawX, drawY
}

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
