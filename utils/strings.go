package utils

import "strings"

const (
	ASCII_ALL         = " !\"#$%'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	ASCII_LETTERS     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ASCII_NON_NUMBERS = "!\"#$%'()*+,-./:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	ASCII_NUMBERS     = "0123456789"
)

func StringDimensions(s string) (int, int) {
	text := strings.Split(s, "\n")
	h := len(text)
	w := 0
	for _, line := range text {
		if len(line) > w {
			w = len(line)
		}
	}

	return w, h
}

func StringDrawLength(s string) int {

	length := 0

	inEscapeCode := false

	for _, ch := range s {
		if ch == '\x1b' && !inEscapeCode {
			inEscapeCode = true
		} else if ch == 'm' && inEscapeCode {
			inEscapeCode = false
		} else if !inEscapeCode {
			length += 1
		}
	}

	return length

}
