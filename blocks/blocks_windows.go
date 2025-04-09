//go:build windows
// +build windows

package blocks

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/mattn/go-isatty"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procSetConsoleCursorPosition   = kernel32.NewProc("SetConsoleCursorPosition")
	procFillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
)

// clear the line and move the cursor up
var clear = fmt.Sprintf("%c[%dA%c[2K\r", ESC, 0, ESC)

type short int16
type dword uint32
type word uint16

type coord struct {
	x short
	y short
}

type smallRect struct {
	left   short
	top    short
	right  short
	bottom short
}

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        word
	window            smallRect
	maximumWindowSize coord
}

func clearLines(out io.Writer, lineCount int) {
	f, ok := out.(FdWriter)
	if ok && !isatty.IsTerminal(f.Fd()) {
		ok = false
	}
	if !ok {
		_, _ = fmt.Fprint(out, strings.Repeat(clear, lineCount))
		return
	}
	fd := f.Fd()
	var csbi consoleScreenBufferInfo
	_, _, _ = procGetConsoleScreenBufferInfo.Call(fd, uintptr(unsafe.Pointer(&csbi)))

	for i := 0; i < lineCount; i++ {
		// move the cursor up
		csbi.cursorPosition.y--
		_, _, _ = procSetConsoleCursorPosition.Call(fd, uintptr(*(*int32)(unsafe.Pointer(&csbi.cursorPosition))))
		// clear the line
		cursor := coord{
			x: csbi.window.left,
			y: csbi.window.top + csbi.cursorPosition.y,
		}
		var count, w dword
		count = dword(csbi.size.x)
		_, _, _ = procFillConsoleOutputCharacter.Call(fd, uintptr(' '), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&w)))
	}
}
