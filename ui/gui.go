package ui

import "github.com/awesome-gocui/gocui"

type GUI struct {
	gui     *gocui.Gui
	widgets map[string]*gocui.Manager
}
