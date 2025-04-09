package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func text(g ttyprogress.Container) {
	for t := 0; t < 2; t++ {
		text, _ := ttyprogress.NewTextSpinner().
			SetPredefined(70).
			SetView(3).
			SetSpeed(1).
			SetFollowUpGap("> ").
			PrependFunc(ttyprogress.Message(fmt.Sprintf("working on task %d...", t+1))).
			AppendElapsed().
			Add(g)

		go func() {
			m := 3 + rand.Int()%8
			for i := 0; i <= m; i++ {
				fmt.Fprintf(text, "doing step %d of task %d\n", i+1, t+1)
				time.Sleep(time.Millisecond * 100 * time.Duration(1+rand.Int()%20))
			}
			text.Close()
		}()
	}
}

func main() {
	p := ttyprogress.For(os.Stdout)
	g, _ := ttyprogress.NewAnonymousGroup().
		HideOnClose().
		SetGap("- ").
		SetFollowUpGap("  ").Add(p)
	text(g)
	g.Close()
	p.Close()
	p.Wait(nil)
}
