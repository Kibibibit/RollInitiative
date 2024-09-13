package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

const DRAW_TABLE_START = "&&DRAW_TABLE_START"

func DrawText(view *gocui.View, colW int, maxH int, colors *ColorPalette, text string, drawX, drawY int) (int, int) {
	view.SetWritePos(drawX, drawY)

	//Contains every line that will eventually be drawn to the screen
	lines := []string{}

	// Break up the string on new lines, as we want to preserve paragraph breaks
	baseLines := strings.Split(text, "\n")
	var index int = 0
	for index < len(baseLines) {
		baseLine := baseLines[index]

		//This means there is a table
		if strings.Contains(baseLine, "|") {
			tableLines := []string{}
			tableLine := baseLines[index]
			for strings.Contains(tableLine, "|") {
				tableLines = append(tableLines, tableLine)
				index++
				if index >= len(baseLines) {
					break
				} else {
					tableLine = baseLines[index]
				}
			}

			table := [][]string{}
			longestLines := []int{}
			for _, line := range tableLines {
				row := []string{}
				colNum := 0
				if strings.Contains(line, "---") {
					continue
				}
				for _, data := range strings.Split(line, "|") {

					cell := strings.TrimSpace(data)
					if len(cell) == 0 {
						continue
					}
					if len(longestLines)-1 < colNum {
						longestLines = append(longestLines, 0)
					}
					if len(cell) > longestLines[colNum] {
						longestLines[colNum] = len(cell)
					}
					row = append(row, cell)
					colNum += 1

				}
				table = append(table, row)

			}

			firstRowChars := []string{}

			for _, l := range longestLines {
				firstRowChars = append(firstRowChars, strings.Repeat("─", l+2))
			}
			drawTableLines := []string{}
			drawTableLines = append(drawTableLines, fmt.Sprintf("┌%s┐", strings.Join(firstRowChars, "┬")))

			for index, row := range table {
				drawItems := []string{}
				for x, cell := range row {
					for len(cell) < longestLines[x] {
						cell = fmt.Sprintf("%s ", cell)
					}
					if index == 0 {
						cell = ApplyBold(cell, colors.FgColor)
					}
					drawItems = append(drawItems, fmt.Sprintf(" %s ", cell))

				}
				drawTableLines = append(drawTableLines, fmt.Sprintf("│%s│", strings.Join(drawItems, "│")))
				if index == 0 {
					drawTableLines = append(drawTableLines, fmt.Sprintf("├%s┤", strings.Join(firstRowChars, "┼")))
				}
			}

			drawTableLines = append(drawTableLines, fmt.Sprintf("└%s┘", strings.Join(firstRowChars, "┴")))

			lines = append(lines, fmt.Sprintf("%s:%d", DRAW_TABLE_START, len(drawTableLines)))
			lines = append(lines, drawTableLines...)

		} else if utils.StringDrawLength(baseLine) < colW-2 {

			if strings.Contains(baseLine, "#") {
				baseLine = strings.ReplaceAll(baseLine, "#", "")
				baseLine = strings.TrimSpace(baseLine)
				baseLine = ApplyBold(baseLine, colors.FgColor)
			}
			lines = append(lines, baseLine)
			index++
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
			index++

		}

	}

	for _, line := range lines {

		view.SetWritePos(drawX, drawY)
		if strings.Contains(line, DRAW_TABLE_START) {
			strAmount := strings.Split(line, ":")[1]
			int64Amount, err := strconv.ParseInt(strAmount, 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			intAmount := int(int64Amount)
			if drawY+intAmount > maxH-1 {
				drawY = 1
				drawX += colW
			}
		} else {
			fmt.Fprint(view, line)
			drawY += 1
			if drawY > maxH-1 {
				drawY = 1
				drawX += colW
			}
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

func RenderPartySeachRow(p *models.Party, colors *ColorPalette, index int, width int) string {
	line := ApplyFgColor(fmt.Sprintf(" %2d ", index), colors.FgColorDim)
	line = fmt.Sprintf(" %s%s", line, p.Name)
	sCount := width - utils.StringDrawLength(line) - 4

	line += strings.Repeat(" ", sCount)

	line += string(len(p.Players))

	for utils.StringDrawLength(line) <= width+2 {
		line += " "
	}

	return line
}
