package main

import (
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	// --- begin bar ---
	bar, _ := ttyprogress.NewSteps("downloading", "unpacking", "installing", "verifying").
		PrependStep().
		PrependFunc(ttyprogress.Message("progressbar"), 0).
		PrependElapsed().AppendCompleted().
		Add(p)
	// --- end bar ---

	bar.Start()
	for i := 0; i < 4; i++ {
		time.Sleep(time.Second * 2)
		bar.Incr()
	}
}
