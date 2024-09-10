package ui

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

var attrList = map[gocui.Attribute]string{
	gocui.AttrBold:      "1",
	gocui.AttrDim:       "2",
	gocui.AttrItalic:    "3",
	gocui.AttrUnderline: "4",
	gocui.AttrBlink:     "5",
	gocui.AttrReverse:   "7",
}

const (
	fgLayer = "38"
	bgLayer = "48"
)

func ApplyStyles(s string, attr gocui.Attribute) string {
	list := []string{}

	for key, sequence := range attrList {
		if key&attr > 0 {
			list = append(list, sequence)
		}
	}

	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", strings.Join(list, ";"), s)

}

func applyColor(s string, layer string, color *Color) string {
	return fmt.Sprintf("\x1b[%s;2;%d;%d;%dm%s\x1b[0m", layer, color.r, color.g, color.b, s)
}

func ApplyFgColor(s string, color *Color) string {
	return applyColor(s, fgLayer, color)
}

func ApplyBgColor(s string, color *Color) string {
	return applyColor(s, bgLayer, color)
}

func ApplyBold(s string, fg *Color) string {
	return ApplyFgColor(ApplyStyles(s, gocui.AttrBold), fg)
}
