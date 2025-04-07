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

	for i, b := range bars {
		bar, _ := ttyprogress.NewSpinner().
			SetPredefined(b).
			SetSpeed(1).
			PrependFunc(ttyprogress.Message(fmt.Sprintf("working on task %d ...", i+1))).
			AppendElapsed().Add(p)
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(10+rand.Int()%20))
			bar.Close()
		}()
	}

	p.Close()
	p.Wait(nil)
}
