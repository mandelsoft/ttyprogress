package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func group(n int, c ttyprogress.Container) {
	anon, _ := ttyprogress.NewAnonymousGroup().Add(c)
	text, _ := ttyprogress.NewText().
		SetTitleLine(fmt.Sprintf("group %d", n)).
		SetGap("# ").
		SetFollowUpGap("> ").
		SetView(3).
		SetAuto().Add(anon)

	if n > 1 {
		group(n-1, anon)
	}
	anon.Close()
	go func() {
		for i := 0; i <= 20; i++ {
			fmt.Fprintf(text, "doing step %d\n", i)
			time.Sleep(time.Millisecond * 500)
		}
		text.Close()
	}()
}

func main() {
	p := ttyprogress.For(os.Stdout)

	group(2, p)
	p.Close()
	p.Wait(nil)
}
