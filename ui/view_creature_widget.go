package ui

import (
	"fmt"
	"log"
	"strings"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

var statNames = []string{
	"STR",
	"DEX",
	"CON",
	"INT",
	"WIS",
	"CHA",
}

const (
	statStringFirst  string = "┌── %s ──"
	statStringMiddle string = "┬── %s ──"
	statStringEnd    string = "┬── %s ──┐"
)

type ViewCreatureWidget struct {
	view           *gocui.View
	name           string
	x, y           int
	w, h           int
	colW           int
	colors         *ColorPalette
	dataStore      *models.DataStore
	previousWidget string
	creature       *models.Creature
	spellX         int
	spellY         int
	hasSpells      bool
	currentSpell   string
}

func NewViewCreatureWidget(dataStore *models.DataStore, previousWidget string, colors *ColorPalette, creature *models.Creature) *ViewCreatureWidget {
	return &ViewCreatureWidget{
		name:           NameViewCreatureWidget,
		dataStore:      dataStore,
		colors:         colors,
		previousWidget: previousWidget,
		creature:       creature,
		spellX:         0,
		spellY:         0,
		hasSpells:      len(creature.Spells) > 0,
		currentSpell:   "",
	}

}

func (w *ViewCreatureWidget) Layout(g *gocui.Gui) error {

	var affinityNames = []string{
		"Damage Immunities",
		"Damage Resistances",
		"Damage Vulnerabilities",
		"Condition Immunities",
		"Senses",
		"Languages",
	}

	width, height := g.Size()

	w.w = utils.Clamp(width-4, 70, width-4)
	w.h = height - 5

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

	w.colW = w.w / 3

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+2, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		w.createKeybinds(g)
		view.Frame = true
		view.BgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.SelBgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.Title = "View Creature"
		view.Visible = true

	} else if err != nil {
		return err
	}

	//Draw Metadata

	view.SetWritePos(1, 1)

	fmt.Fprint(view,
		ApplyBold(
			fmt.Sprintf("%s - CR %s (%d XP, +%d PB)", w.creature.Name, w.creature.CR, models.XPFromCR(w.creature.CR), models.PBFromCR(w.creature.CR)),
			w.colors.FgColor,
		),
	)
	view.SetWritePos(1, 2)
	fmt.Fprintf(view, "- %s, %s, %s (%s)", w.creature.Type, w.creature.Size, w.creature.Alignment, w.creature.Source)

	view.SetWritePos(1, 4)
	fmt.Fprintf(view, "%s %s", ApplyBold("Armour Class:", w.colors.FgColor), w.creature.AC)

	conMod := w.creature.GetConMod()

	hpBoost := conMod * w.creature.HitDice

	view.SetWritePos(1, 5)
	hpString := fmt.Sprintf("%s %d (%dd%d + %d)", ApplyBold("Hit Points:", w.colors.FgColor), utils.AverageDiceRoll(w.creature.HitDice, w.creature.HitDiceType)+hpBoost, w.creature.HitDice, w.creature.HitDiceType, conMod*w.creature.HitDice)
	hpString = strings.ReplaceAll(hpString, " + -", " - ")
	hpString = strings.ReplaceAll(hpString, " + 0", "")
	fmt.Fprint(view, hpString)

	view.SetWritePos(1, 6)
	fmt.Fprintf(view, "%s %s", ApplyBold("Speed:", w.colors.FgColor), w.creature.Speed)

	//Draw statblock

	topLine := ""

	bottomLine := "└─────────┴─────────┴─────────┴─────────┴─────────┴─────────┘"

	statLine := ""
	for index, value := range w.creature.Stats {
		statName := statNames[index]
		statName = ApplyBold(statName, w.colors.FgColor)
		statString := statStringMiddle
		if index == 0 {
			statString = statStringFirst
		} else if index == len(statNames)-1 {
			statString = statStringEnd
		}
		topLine += fmt.Sprintf(statString, statName)

		statLine += fmt.Sprintf("│ %2d (%2d) ", value, w.creature.GetMod(value))
	}
	statLine += "│"

	statLine = strings.ReplaceAll(statLine, "( ", "(+")
	view.SetWritePos(1, 8)
	fmt.Fprint(view, topLine)
	view.SetWritePos(1, 9)
	fmt.Fprint(view, statLine)
	view.SetWritePos(1, 10)
	fmt.Fprint(view, bottomLine)

	var drawX, drawY int
	drawX = 1
	drawY = 12

	view.SetWritePos(drawX, drawY)

	// Draw saves and skills
	view.SetWritePos(drawX, drawY)

	saves, hasSaves := w.getSkillSavesString("Saving Throws", w.creature.Saves)
	if hasSaves {
		fmt.Fprint(view, saves)
		drawY += 1
	}

	view.SetWritePos(drawX, drawY)

	skills, hasSkills := w.getSkillSavesString("Skills", w.creature.Skills)
	if hasSkills {
		fmt.Fprint(view, skills)
		drawY += 1
	}

	//resistences, senses, languages and so on
	var affinities = []string{
		w.creature.DamageImmunities,
		w.creature.DamageResistances,
		w.creature.DamageVulnerabilities,
		w.creature.ConditionImmunities,
		w.creature.Senses,
		w.creature.Languages,
	}

	for index, affinity := range affinities {

		if len(affinity) > 0 {
			drawX, drawY = w.drawText(fmt.Sprintf("%s: %s", ApplyBold(affinityNames[index], w.colors.FgColor), affinity), drawX, drawY)
		}
	}

	drawX, drawY = w.drawCreatureTraitList("", w.creature.Traits, drawX, drawY)

	drawX, drawY = w.drawCreatureSpellList(drawX, drawY)

	drawX, drawY = w.drawCreatureTraitList("Actions", w.creature.Actions, drawX, drawY)
	drawX, drawY = w.drawCreatureTraitList("Bonus Actions", w.creature.BonusActions, drawX, drawY)
	drawX, drawY = w.drawCreatureTraitList("Reactions", w.creature.Reactions, drawX, drawY)
	drawX, drawY = w.drawCreatureTraitList("Lair Actions", w.creature.LairActions, drawX, drawY)
	//TODO: Legendary descriptions
	w.drawCreatureTraitList("Legendary Actions", w.creature.LegendaryActions, drawX, drawY)

	return nil

}

func (w *ViewCreatureWidget) drawText(text string, drawX, drawY int) (int, int) {
	return DrawText(w.view, w.colW, w.h, w.colors, text, drawX, drawY)
}

func (w *ViewCreatureWidget) drawCreatureSpellList(drawX, drawY int) (int, int) {
	var spellLevelTitles = []string{
		"Cantrips",
		"1st level",
		"2nd level",
		"3rd level",
		"4th level",
		"5th level",
		"6th level",
		"7th level",
		"8th level",
		"9th level",
	}
	if len(w.creature.Spells) > 0 {
		drawY += 1
		drawX, drawY = w.drawText(fmt.Sprintf("%s: %s", ApplyBold("Spellcasting", w.colors.FgColor), w.creature.SpellNotes), drawX, drawY)
		for level := 0; level <= 9; level++ {
			if spells, ok := w.creature.Spells[level]; ok {
				drawLine := spellLevelTitles[level]
				slotsString := "(at will)"
				if level > 0 {
					slotsString = fmt.Sprintf("(%d slots)", spells.Slots)
				}

				drawLine = fmt.Sprintf("%s %s:", drawLine, slotsString)

				spellNames := []string{}

				for x, spell := range spells.Spells {
					s := w.dataStore.GetSpell(spell)

					spellDrawString := spell
					if s != nil {
						spellDrawString = s.Name
					}

					if level == w.spellY && x == w.spellX {
						spellDrawString = fmt.Sprintf("\x1b[7m%s\x1b[0m", spellDrawString)
						w.currentSpell = s.Id
					}

					spellNames = append(spellNames, spellDrawString)
				}

				drawLine = fmt.Sprintf("%s %s", drawLine, strings.Join(spellNames, ", "))

				drawX, drawY = w.drawText(drawLine, drawX, drawY)
			}
		}
		if len(w.creature.PrecombatSpells) > 0 {
			drawLine := "The creature casts the following spells on itself before combat:"
			spellNames := []string{}
			for _, spell := range w.creature.PrecombatSpells {
				s := w.dataStore.GetSpell(spell)
				if s == nil {
					spellNames = append(spellNames, spell)
				} else {
					spellNames = append(spellNames, s.Name)
				}
			}
			drawLine = fmt.Sprintf("%s %s", drawLine, strings.Join(spellNames, ", "))
			drawX, drawY = w.drawText(drawLine, drawX, drawY)
		}

	}
	return drawX, drawY
}

func (w *ViewCreatureWidget) drawCreatureTraitList(title string, list []models.CreatureTrait, drawX, drawY int) (int, int) {
	if len(list) > 0 {

		drawY += 1

		w.view.SetWritePos(drawX, drawY)
		if len(title) > 0 {
			fmt.Fprint(w.view, ApplyBold(ApplyStyles(title, gocui.AttrUnderline), w.colors.FgColor))
			drawY += 1
		}

		for _, trait := range list {
			drawX, drawY = w.drawText(fmt.Sprintf("%s: %s", ApplyBold(trait.Name, w.colors.FgColor), trait.Description), drawX, drawY)
			drawY += 1
		}
	}

	return drawX, drawY
}

func (w *ViewCreatureWidget) getSkillSavesString(title string, data map[string]int) (line string, hasSaves bool) {
	skillsaves := []string{}
	for stat, bonus := range data {
		skillsaves = append(skillsaves, strings.ReplaceAll(fmt.Sprintf("%s +%d", stat, bonus), "+-", "-"))
	}
	if len(skillsaves) == 0 {
		return "", false
	}
	out := ApplyBold(fmt.Sprintf("%s:", title), w.colors.FgColor)
	out = fmt.Sprintf("%s %s", out, strings.Join(skillsaves, ", "))

	return out, true

}

func (w *ViewCreatureWidget) createKeybinds(g *gocui.Gui) error {

	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.onClose); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.onArrow(0, 1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.onArrow(0, -1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.onArrow(-1, 0)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.onArrow(1, 0)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.onEnter); err != nil {
		return err
	}

	return nil

}

func (w *ViewCreatureWidget) onEnter(g *gocui.Gui, _ *gocui.View) error {
	if w.hasSpells {
		if len(w.currentSpell) > 0 {
			viewSpell := NewViewSpellWidget(w.dataStore, w.name, w.colors, w.dataStore.GetSpell(w.currentSpell))

			viewSpell.Layout(g)

			g.Update(func(g *gocui.Gui) error {

				_, err := g.SetCurrentView(viewSpell.name)
				return err
			})
		}
	}

	return nil
}

func (w *ViewCreatureWidget) onArrow(xOffset, yOffset int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		if w.hasSpells {

			w.spellY += yOffset

			for len(w.creature.Spells[w.spellY].Spells) == 0 {
				w.spellY += yOffset
				if w.spellY < 0 {
					w.spellY = 9
				}
				if w.spellY > 9 {
					w.spellY = 0

				}
			}

			xSize := len(w.creature.Spells[w.spellY].Spells)

			w.spellX = utils.Clamp(w.spellX, 0, xSize-1)
			w.spellX += xOffset
			if w.spellX < 0 {
				w.spellX += xSize
			}
			if w.spellX >= xSize {
				w.spellX -= xSize
			}

			w.Layout(g)

		}

		return nil
	}
}

func (w *ViewCreatureWidget) onClose(g *gocui.Gui, v *gocui.View) error {

	nextView, err := g.View(NameRootWidget)
	if err != nil {
		log.Println(err)
		return err
	}

	w.view.Visible = false
	g.DeleteKeybindings(w.name)
	g.DeleteView(w.name)
	w.view.Frame = false
	w.view.Clear()

	g.Update(func(g *gocui.Gui) error {

		_, err := g.SetCurrentView(nextView.Name())
		return err
	})

	return nil
}
