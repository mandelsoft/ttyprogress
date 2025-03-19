package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func main() {
	p := ttyprogress.New(os.Stdout)

	for s := 0; s < 3; s++ {
		bar, _ := ttyprogress.NewTextSpinner().
			SetPredefined(5).
			SetView(3).
			SetFollowUpGap("> ").
			PrependFunc(ttyprogress.Message(fmt.Sprintf("working on task %d...", s+1))).
			AppendElapsed().
			Add(p)

		go func() {
			// starts automatically, with the first write
			steps := 6 + rand.Int()%10
			for i := 0; i <= steps; i++ {
				t := 500 + 200*(rand.Int()%6)
				fmt.Fprintf(bar, "doing step %d\n", i)
				time.Sleep(time.Duration(t) * time.Millisecond)
			}
			bar.Close()
		}()
	}
	p.Close()
	p.Wait(nil)
}
