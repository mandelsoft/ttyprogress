package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/synclog"
)

func Step(n string) ttyprogress.NestedStep {
	return ttyprogress.NewNestedStep[ttyprogress.Bar](
		n, ttyprogress.NewBar().SetTotal(100).
			// HideOnClose().
			PrependElapsed().
			AppendCompleted(),
	)
}

func main() {
	synclog.LogToFile("local/lock.log")
	p := ttyprogress.For(os.Stdout)

	bar, _ := ttyprogress.NewNestedSteps(
		Step("downloading"),
		Step("unpacking"),
		Step("installing"),
		Step("verifying")).
		SetGap("  ").
		SetWidth(40).
		ShowStepTitle(false).
		PrependFunc(ttyprogress.Message("progressbar"), 0).
		PrependElapsed().
		AppendCompleted().
		Add(p)

	p.Close()

	go func() {
		bar.Start()
		e := bar.Current()
		for i := 0; i < 4; i++ {
			for i := 0; i < 100; i++ {
				time.Sleep(time.Millisecond * time.Duration(rand.Int()%100))
				e.(ttyprogress.Bar).Incr()
			}
			e, _ = bar.Incr()
		}
	}()

	p.Wait(nil)
}
