package ui

import (
	"fmt"
	"strings"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

// This widget will get a string

type StringInputWidget struct {
	view           *gocui.View
	name           string
	title          string
	colors         *ColorPalette
	x, y           int
	w, h           int
	onSubmit       func(v string)
	charSet        string
	previousWidget string
	killed         bool
	data           string
	cursorX        int
}

func NewStringInputWidget(name string, title string, colors *ColorPalette, previousWidget string, charSet string, data string, onSubmit func(v string)) *StringInputWidget {
	return &StringInputWidget{
		name:           name,
		title:          title,
		onSubmit:       onSubmit,
		colors:         colors,
		killed:         false,
		charSet:        charSet,
		data:           data,
		previousWidget: previousWidget,
		cursorX:        len(data),
	}
}

func (w *StringInputWidget) Layout(g *gocui.Gui) error {
	if w.killed {
		g.DeleteKeybindings(w.name)
		g.DeleteView(w.name)
		w.view.Clear()
		w.view.Frame = false

		return nil
	}

	width, height := g.Size()

	w.w = width / 4
	w.h = 1

	w.x = width/2 - w.w/2 - 1
	w.y = height/2 - w.h/2 - 1

	view, err := g.SetView(w.name, w.x, w.y, w.x+w.w+2, w.y+w.h+1, 0)
	w.view = view

	if err == gocui.ErrUnknownView {
		view.BgColor = w.colors.BgColorWindow.GetCUIAttr()
		view.FgColor = w.colors.FgColor.GetCUIAttr()
		view.Title = w.title
		view.SelBgColor = w.colors.FgColor.GetCUIAttr()
		view.SelFgColor = w.colors.BgColor.GetCUIAttr()
		w.setKeybinding(g)
	} else if err != nil {
		return err
	}

	w.view.Clear()

	drawString := []rune{}

	for i, ch := range w.data {
		if i != w.cursorX {
			drawString = append(drawString, ch)
		} else {
			drawString = append(drawString, []rune(w.drawCursor(string(ch)))...)
		}
	}
	if w.cursorX == len(w.data) {
		drawString = append(drawString, []rune(w.drawCursor(" "))...)
	}

	fmt.Fprint(w.view, string(drawString))

	return nil
}

func (w *StringInputWidget) drawCursor(line string) string {
	return fmt.Sprintf("\x1b[7m%s\x1b[0m", line)
}

func (w *StringInputWidget) setKeybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.Kill); err != nil {
		return err
	}

	for _, ch := range w.charSet {
		if err := g.SetKeybinding(w.name, ch, gocui.ModNone, w.onLetter(ch)); err != nil {
			return err
		}
	}
	if strings.Contains(w.charSet, " ") {
		if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.onLetter(' ')); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding(w.name, gocui.KeyBackspace|gocui.KeyBackspace2, gocui.ModNone, w.onDeleteChar(1, 1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyDelete, gocui.ModNone, w.onDeleteChar(0, 0)); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.onMoveCursor(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.onMoveCursor(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModShift, w.onMoveCursor(-len(w.data))); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModShift, w.onMoveCursor(len(w.data))); err != nil {
		return err
	}

	if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.onEnter); err != nil {
		return err
	}

	return nil
}

func (w *StringInputWidget) Kill(g *gocui.Gui, v *gocui.View) error {
	w.killed = true

	g.Update(
		func(g *gocui.Gui) error {
			w.Layout(g)
			_, err := g.SetCurrentView(w.previousWidget)
			return err
		})
	return nil
}

func (w *StringInputWidget) onLetter(ch rune) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		partA := w.data[0:w.cursorX]
		partB := w.data[w.cursorX:len(w.data)]

		w.data = fmt.Sprintf("%s%c%s", partA, ch, partB)

		w.cursorX += 1

		return w.Layout(g)
	}
}

func (w *StringInputWidget) onDeleteChar(charOffset int, cursorOffset int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if len(w.data) > 0 {

			newData := []rune{}
			for i, ch := range w.data {
				if w.cursorX-charOffset >= 0 {
					if i != w.cursorX-charOffset {
						newData = append(newData, ch)
					}
				}
			}
			w.data = string(newData)
			w.cursorX -= cursorOffset
		}
		return w.Layout(g)
	}
}

func (w *StringInputWidget) onEnter(g *gocui.Gui, v *gocui.View) error {
	w.Kill(g, v)
	w.onSubmit(w.data)
	return nil
}

func (w *StringInputWidget) onMoveCursor(offset int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		w.cursorX = utils.Clamp(w.cursorX+offset, 0, len(w.data))
		return w.Layout(g)
	}
}
