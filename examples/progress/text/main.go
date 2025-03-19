package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func main() {
	p := ttyprogress.New(os.Stdout)

	text, _ := ttyprogress.NewText().
		SetTitleLine("some output").
		SetFollowUpGap("> ").
		SetView(3).
		SetAuto().Add(p)

	for i := 0; i <= 20; i++ {
		fmt.Fprintf(text, "doing step %d\n", i)
		time.Sleep(time.Millisecond * 500)
	}
	text.Close()
}
