package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"windmills/roll_initiative/models"
	"windmills/roll_initiative/ui"
)

func main() {

	defaultFolderString := os.Getenv("ROLL_DATA_FOLDERS")

	if len(defaultFolderString) == 0 {
		log.Println("Warning! ROLL_DATA_FOLDER not set, defaulting to ./srd_data for data location!")
		defaultFolderString = "./srd_data"
	}

	dataPathPtr := flag.String("data", defaultFolderString, "The location of the data folder")

	flag.Parse()

	dataPathStr := *dataPathPtr

	dataFolders := strings.Split(dataPathStr, ",")

	dataStore := models.MakeDataStore(dataFolders)

	colors := ui.ColorPalette{
		BgColor:       ui.NewColor(29, 31, 48),
		FgColor:       ui.NewColor(255, 255, 255),
		BgColorWindow: ui.NewColor(33, 35, 54),
		FgColorDim:    ui.NewColor(103, 105, 118),
	}

	gui, err := ui.CreateGui(dataStore, &colors)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer gui.Close()

	ui.MainLoop(gui)

}
