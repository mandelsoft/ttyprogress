package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func grouped(p ttyprogress.Container, lvl int) {
	s := ttyprogress.NewSpinner().SetPredefined(86).
		SetSpeed(5).
		PrependFunc(ttyprogress.Message(fmt.Sprintf("Grouped work %d", lvl))).
		AppendElapsed()

	g, _ := ttyprogress.NewGroup[ttyprogress.Spinner](s).
		SetGap("- ").
		SetFollowUpGap("  ").
		Add(p)
	if lvl > 0 {
		grouped(g, lvl-1)
	}
	for i := 0; i < 2; i++ {
		bar, _ := ttyprogress.NewSpinner().SetPredefined(70).
			SetSpeed(1).
			PrependFunc(ttyprogress.Message(fmt.Sprintf("working on task %d[%d]...", i+1, lvl))).
			AppendElapsed().
			Add(g)
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(10+rand.Int()%20))
			bar.Close()
		}()
	}

	text, _ := ttyprogress.NewTextSpinner().
		SetPredefined(70).
		SetView(3).
		SetSpeed(1).
		SetFollowUpGap("  ").
		PrependFunc(ttyprogress.Message(fmt.Sprintf("working on task %d[%d]...", 3, lvl))).
		AppendElapsed().
		Add(g)

	go func() {
		for i := 0; i <= 20; i++ {
			fmt.Fprintf(text, "doing step %d of task %d[%d]\n", i, 3, lvl)
			time.Sleep(time.Millisecond * 100 * time.Duration(1+rand.Int()%20))
		}
		text.Close()
	}()
	g.Close()
}

func main() {
	p := ttyprogress.New(os.Stdout)
	grouped(p, 2)
	p.Close()
	p.Wait(nil)
}
