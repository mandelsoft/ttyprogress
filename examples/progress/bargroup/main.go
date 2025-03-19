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

func text(g ttyprogress.Group) {

	text, _ := ttyprogress.NewTextSpinner().
		SetPredefined(70).
		SetView(3).
		SetSpeed(1).
		SetFollowUpGap("> ").
		PrependFunc(ttyprogress.Message(fmt.Sprintf("working on task %d...", 3))).
		AppendElapsed().
		Add(g)

	go func() {
		for i := 0; i <= 20; i++ {
			fmt.Fprintf(text, "doing step %d of task %d\n", i, 3)
			time.Sleep(time.Millisecond * 100 * time.Duration(1+rand.Int()%20))
		}
		text.Close()
	}()
}

func main() {
	p := ttyprogress.New(os.Stdout)
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
