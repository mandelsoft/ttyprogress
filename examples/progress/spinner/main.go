package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/mandelsoft/ttyprogress"
)

func main() {
	p := ttyprogress.New(os.Stdout)

	bars := []int{1000, 1004, 1003, 1002}
	cols := []*color.Color{
		color.New(color.FgHiYellow, color.Bold),
		color.New(color.FgCyan, color.Italic),
		color.New(color.BgGreen, color.Underline),
		color.New(color.FgGreen),
	}
	for i, b := range bars {
		bar, _ := ttyprogress.NewSpinner().
			SetPredefined(b).
			SetSpeed(1).
			SetColor(cols[i]).
			PrependFunc(ttyprogress.Message(fmt.Sprintf("very long line working on %d ...", i+1))).
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
