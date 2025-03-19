package blocks

import (
	"github.com/fatih/color"
)

var _ = color.Reset

// ColorLength works together with color.New.
// and returns the byte length for a color format escape sequence
// if the data starts with such a prefix.
func ColorLength(data []byte) int {
	if len(data) <= 3 || data[0] != '\x1b' || data[1] != '[' {
		return 0
	}

	for i, c := range data[2:] {
		if (c < '0' || c > '9') && c != ';' {
			if c != 'm' {
				return 0
			}
			return i + 3
		}
	}
	return 0
}
