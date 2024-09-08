package utils

import "strings"

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
