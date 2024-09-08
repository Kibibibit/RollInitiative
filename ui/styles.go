package ui

import (
	"fmt"
	"windmills/roll_initiative/utils"

	"github.com/awesome-gocui/gocui"
)

type Color struct {
	r int
	g int
	b int
}

type ColorPalette struct {
	BGColor       *Color
	FGColor       *Color
	WindowBGColor *Color
}

func (c *Color) GetCUIAttr() gocui.Attribute {
	return gocui.NewRGBColor(int32(c.r), int32(c.g), int32(c.b))
}

func NewColor(r int, g int, b int) *Color {
	return &Color{
		r: utils.Clamp(r, 0, 255),
		g: utils.Clamp(g, 0, 255),
		b: utils.Clamp(b, 0, 255),
	}
}

func FprintRGB(v *gocui.View, c *Color, data string) (int, error) {
	return fmt.Fprintf(
		v,
		"\x1b[38;2;%d;%d;%dm%s\x1b[0m",
		c.r, c.g, c.b, data,
	)
}
