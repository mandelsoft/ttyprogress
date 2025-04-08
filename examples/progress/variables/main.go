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
		PrependVariable("temperature").
		PrependMessage("...")

	for i, b := range bars {
		bar, _ := def.SetPredefined(b).Add(p)
		bar.SetVariable("name", fmt.Sprintf("task %d", i+1))
		bar.Start()
		go func() {
			temp := 10
			for i := 0; i < 10; i++ {
				temp += rand.IntN(5) - 2
				bar.SetVariable("temperature", fmt.Sprintf("[%d°]", temp))
				time.Sleep(time.Second * time.Duration(rand.IntN(3)))
			}
			bar.SetVariable("temperature", "")
			bar.Close()
		}()
	}

	p.Close()
	p.Wait(nil)
}
