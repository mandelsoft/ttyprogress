package main

import (
	"context"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/units"
)

func main() {
	// setup output context
	p := ttyprogress.For(os.Stdout)

	// configure indicator
	bar := ttyprogress.NewBar().
		SetPredefined(10).
		SetTotal(500).
		SetWidth(ttyprogress.PercentTerminalSize(30)).
		PrependMessage("Downloading...").
		PrependElapsed().AppendCompleted().
		AppendFunc(ttyprogress.Amount(units.Bytes(1024)))

	ttyprogress.RunWith(p, bar, func(bar ttyprogress.Bar) {
		for i := 0; i <= 20; i++ {
			bar.Set(i * 5 * 5)
			time.Sleep(time.Millisecond * 500)
		}
	})

	p.Close()
	p.Wait(context.Background())
}
