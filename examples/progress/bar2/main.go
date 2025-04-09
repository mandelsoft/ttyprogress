package main

import (
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/units"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	bar1, _ := ttyprogress.NewBar().
		SetWidth(40).
		SetPredefined(10).
		SetTotal(500).
		PrependMessage("Downloading...").
		AppendFunc(ttyprogress.Amount(units.Bytes(1024))).
		Add(p)

	bar2, _ := ttyprogress.NewBar().
		SetWidth(40).
		SetPredefined(10).
		SetTotal(500).
		PrependMessage("Downloading...").
		PrependElapsed().AppendCompleted().
		AppendFunc(ttyprogress.Amount(units.Bytes(1024))).
		Add(p)

	for i := 0; i <= 20; i++ {
		bar2.Set(i * 5 * 5)
		if i%2 == 0 {
			bar1.Set(i * 5 * 5)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
