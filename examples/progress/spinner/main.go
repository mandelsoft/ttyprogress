package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	bars := []int{1000, 1002, 1003}
	cols := []ttycolors.Format{
		ttycolors.New(ttycolors.FmtBrightGreen, ttycolors.FmtUnderline),
		ttycolors.New(ttycolors.FmtCyan, ttycolors.FmtItalic),
		ttycolors.New(ttycolors.FmtBgCyan, ttycolors.FmtBold),
	}
	for i, b := range bars {
		bar, _ := ttyprogress.NewSpinner().
			SetPredefined(b).
			SetSpeed(1).
			SetColor(cols[i]).
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
