package ui

import (
	"errors"
	"log"
	"windmills/roll_initiative/models"

	"github.com/awesome-gocui/gocui"
)

func CreateGui(dataStore *models.DataStore, colors *ColorPalette) (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.OutputTrue, true)
	if err != nil {
		log.Fatalln("Failed to create new GUI!")
		return nil, err
	}

	g.BgColor = colors.BGColor.GetCUIAttr()
	g.FgColor = colors.FGColor.GetCUIAttr()

	rootWidget := NewRootWidget(dataStore, colors)
	shortcutsWidget := NewShortcutsWidget(dataStore, colors)

	g.SetManager(rootWidget, shortcutsWidget)

	g.Update(func(g *gocui.Gui) error {
		g.SetCurrentView(rootWidget.name)
		return nil
	})

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		log.Panicln("Failed to set ctrlC keybinding!")
		return nil, err
	}

	return g, nil
}

func MainLoop(g *gocui.Gui) {
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}
