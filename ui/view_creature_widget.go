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

var affinityNames = []string{
	"Damage Immunities",
	"Damage Resistances",
	"Damage Vulnerabilities",
	"Condition Immunities",
	"Senses",
	"Languages",
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
	colors         *ColorPalette
	dataStore      *models.DataStore
	previousWidget string
	creature       *models.Creature
}

func NewViewCreatureWidget(dataStore *models.DataStore, previousWidget string, colors *ColorPalette, creature *models.Creature) *ViewCreatureWidget {
	return &ViewCreatureWidget{
		name:           NameViewCreatureWidget,
		dataStore:      dataStore,
		colors:         colors,
		previousWidget: previousWidget,
		creature:       creature,
	}

}

func (w *ViewCreatureWidget) Layout(g *gocui.Gui) error {

	width, height := g.Size()

	w.w = utils.Clamp(width-10, height-5, width)
	w.h = height - 10

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

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
			fmt.Sprintf("%s - CR %s (%d XP)", w.creature.Name, w.creature.CR, models.XPFromCR(w.creature.CR)),
			w.colors.FgColor,
		),
	)
	view.SetWritePos(1, 2)
	fmt.Fprintf(view, "- %s, %s, %s (%s)", w.creature.Type, w.creature.Size, w.creature.Alignment, w.creature.Source)

	view.SetWritePos(1, 4)
	fmt.Fprintf(view, "%s %s", ApplyBold("Armour Class:", w.colors.FgColor), w.creature.AC)

	conMod := (w.creature.Stats[2] - 10) / 2

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

		statLine += fmt.Sprintf("│ %2d (%2d) ", value, ((value - 10) / 2))
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
		view.SetWritePos(drawX, drawY)

		if len(affinity) > 0 {
			fmt.Fprintf(view, "%s: %s", ApplyBold(affinityNames[index], w.colors.FgColor), affinity)
			drawY += 1
		}
	}

	drawY += 1

	//Draw traits

	for _, trait := range w.creature.Traits {
		drawX, drawY = w.drawCreatureTrait(&trait, drawX, drawY)
	}

	drawY += 1

	if len(w.creature.Actions) > 0 {
		view.SetWritePos(drawX, drawY)
		fmt.Fprint(view, ApplyBold(ApplyStyles("Actions", gocui.AttrUnderline), w.colors.FgColor))
		drawY += 1

		for _, action := range w.creature.Actions {
			drawX, drawY = w.drawCreatureTrait(&action, drawX, drawY)
		}
	}
	if len(w.creature.LegendaryActions) > 0 {

		drawY += 1

		view.SetWritePos(drawX, drawY)
		fmt.Fprint(view, ApplyBold(ApplyStyles("Legendary Actions", gocui.AttrUnderline), w.colors.FgColor))
		drawY += 1

		for _, action := range w.creature.LegendaryActions {
			drawX, drawY = w.drawCreatureTrait(&action, drawX, drawY)
		}
	}

	return nil

}

func (w *ViewCreatureWidget) drawCreatureTrait(trait *models.CreatureTrait, drawX int, drawY int) (int, int) {
	drawYOut := drawY
	drawXOut := drawX
	w.view.SetWritePos(drawX, drawY)

	drawLine := trait.Description
	lineLength := len(strings.Split(trait.Description, "\n"))
	if lineLength > 0 {
		drawLine = strings.ReplaceAll(drawLine, "\n", "\n\t")
	}
	fmt.Fprintf(w.view, "%s: %s", ApplyBold(trait.Name, w.colors.FgColor), drawLine)

	drawYOut += lineLength

	return drawXOut, drawYOut
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
	return nil

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
