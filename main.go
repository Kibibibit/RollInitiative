package main

import (
	"log"
	"os"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/ui"
)

func main() {

	dataStore := models.MakeDataStore()

	colors := ui.ColorPalette{
		BgColor:       ui.NewColor(29, 31, 48),
		FgColor:       ui.NewColor(255, 255, 255),
		BgColorWindow: ui.NewColor(33, 35, 54),
		FgColorDim:    ui.NewColor(103, 105, 118),
	}

	gui, err := ui.CreateGui(dataStore, &colors)

	defer gui.Close()

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	ui.MainLoop(gui)

}
