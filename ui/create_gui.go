package ui

import (
	"errors"
	"log"
	"os"
	"windmills/roll_initiative/models"

	"github.com/awesome-gocui/gocui"
)

func CreateGui(dataStore *models.DataStore, colors *ColorPalette) (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.OutputTrue, true)
	if err != nil {
		log.Fatalln("Failed to create new GUI!")
		return nil, err
	}

	g.BgColor = colors.BgColor.GetCUIAttr()
	g.FgColor = colors.FgColor.GetCUIAttr()

	rootWidget := NewRootWidget(dataStore, colors)
	shortcutsWidget := NewShortcutsWidget(rootWidget, dataStore, colors)

	g.SetManager(rootWidget, shortcutsWidget)

	g.Update(func(g *gocui.Gui) error {
		g.SetCurrentView(rootWidget.name)
		return nil
	})

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit(colors)); err != nil {
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

func Quit(colors *ColorPalette) func(g *gocui.Gui, v *gocui.View) error {

	return func(g *gocui.Gui, v *gocui.View) error {

		currentView := g.CurrentView()

		widget := NewConfirmWidget("Really quit?", colors, currentView.Name(), false, "Are you sure you want to quit?", "Yes", "No",

			func(b bool) {
				if b {
					g.Close()
					os.Exit(0)
				}
			},
		)

		widget.Layout(g)

		g.Update(func(g *gocui.Gui) error {
			g.SetCurrentView(widget.name)
			return nil
		})

		return nil
	}

}
