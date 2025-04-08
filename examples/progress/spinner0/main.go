package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/specs"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	bars := []int{1000, specs.SpinnerType}

	def := ttyprogress.NewSpinner().
		SetSpeed(1).
		AppendElapsed().
		PrependMessage("working on").
		PrependVariable("name").
		PrependMessage("...")

	for i, b := range bars {
		bar, _ := def.SetPredefined(b).Add(p)
		bar.SetVariable("name", fmt.Sprintf("task %d", i+1))
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(10+rand.Int()%20))
			bar.Close()
		}()
	}

	p.Close()
	p.Wait(nil)
}
