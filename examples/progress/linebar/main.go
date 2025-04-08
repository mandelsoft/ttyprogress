package main

import (
	"os"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/units"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	bar, _ := ttyprogress.NewLineBar().
		SetTotal(500).
		PrependMessage("Downloading...").
		PrependElapsed().AppendCompleted().SetDecoratorFormat(ttycolors.FmtGreen).
		AppendFunc(ttyprogress.Amount(units.Bytes(1024))).
		Add(p)

	for i := 0; i <= 20; i++ {
		bar.Set(i * 5 * 5)
		time.Sleep(time.Millisecond * 500)
	}
}
