package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	// --- begin text ---
	text, _ := ttyprogress.NewText().
		SetTitleLine("some output").
		SetFollowUpGap("> ").
		SetView(3).
		SetAuto().Add(p)
	// --- end text ---

	for i := 0; i <= 20; i++ {
		fmt.Fprintf(text, "doing step %d\n", i)
		time.Sleep(time.Millisecond * 500)
	}
	text.Close()
}
