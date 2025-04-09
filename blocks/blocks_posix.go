//go:build !windows
// +build !windows

package blocks

import (
	"fmt"
	"io"
	"strings"
)

// clear the line and move the cursor up
var clear = fmt.Sprintf("%c[%dA%c[2K", ESC, 1, ESC)

func clearLines(out io.Writer, lineCount int) {
	_, _ = fmt.Fprint(out, strings.Repeat(clear, lineCount))
}
