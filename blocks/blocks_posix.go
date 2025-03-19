//go:build !windows
// +build !windows

package blocks

import (
	"fmt"
	"strings"
)

// clear the line and move the cursor up
var clear = fmt.Sprintf("%c[%dA%c[2K", ESC, 1, ESC)

func (w *Blocks) clearLines() {
	_, _ = fmt.Fprint(w.out, strings.Repeat(clear, w.lineCount))
}
