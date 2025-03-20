package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func spinner(g ttyprogress.Group) {
	for i := 0; i < 3; i++ {
		bar, _ := ttyprogress.NewSpinner().
			SetPredefined(70).
			SetSpeed(1).
			PrependFunc(ttyprogress.Message(fmt.Sprintf("working on task %d...", i+1))).
			AppendElapsed().
			Add(g)
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(10+rand.Int()%20))
			bar.Close()
		}()
	}
}

func main() {
	p := ttyprogress.For(os.Stdout)

	// use progress bar to indicate group progress.
	s := ttyprogress.NewBar().
		SetPredefined(1).
		SetWidth(30).
		PrependFunc(ttyprogress.Message(fmt.Sprintf("Grouped work"))).
		// AppendElapsed().
		AppendCompleted()

	g, _ := ttyprogress.NewGroup[ttyprogress.Bar](s).
		SetGap("* ").
		SetFollowUpGap("  ").Add(p)
	spinner(g)
	g.Close()
	p.Close()
	p.Wait(nil)
}
