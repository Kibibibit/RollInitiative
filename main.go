package main

import (
	"fmt"
	"os"
	"os/signal"

	gc "github.com/rthornton128/goncurses"
)

func createNewWin(height int, width int, starty int, startx int) *gc.Window {
	var localWindow *gc.Window

	localWindow, err := gc.NewWindow(height, width, starty, startx)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	localWindow.Box(0, 0)

	localWindow.Refresh()

	return localWindow

}

func destroyWin(window *gc.Window) {
	window.Border(' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ')
	window.Refresh()
	window.Delete()
}

func cleanup() {
	gc.Cursor(1)
	gc.End()
}

func main() {

	var window *gc.Window

	var startx, starty, width, height int

	stdscr, err := gc.Init()

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, os.Interrupt)
	go func() {
		<-closeSignal
		cleanup()
		os.Exit(1)
	}()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer gc.Cursor(1)
	defer gc.End()

	LINES, COLS := stdscr.MaxYX()

	gc.CBreak(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)

	height = 3
	width = 10

	starty = (LINES - height) / 2
	startx = (COLS - width) / 2

	stdscr.Print("Press F1 to exit" + string('ƫ'))
	stdscr.Refresh()

	window = createNewWin(height, width, starty, startx)

	var ch gc.Key

	for ch != gc.KEY_F1 {
		ch = stdscr.GetChar()

		switch ch {
		case gc.KEY_LEFT:
			destroyWin(window)
			startx -= 1

		}

		window = createNewWin(height, width, starty, startx)
	}

}
