package main

// --- begin usage ---
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

	// configure indicator (types)
	// --- begin bar ---
	bar := ttyprogress.NewBar().
		SetPredefined(10).
		SetTotal(500).
		SetWidth(ttyprogress.PercentTerminalSize(30)).
		PrependMessage("Downloading...").
		PrependElapsed().AppendCompleted().
		AppendFunc(ttyprogress.Amount(units.Bytes(1024)))
	// --- end bar ---

	// instantiate bar definition for context
	// (similar to elem, _ := bar.Add(p) )
	// and run some process updating the progress on the indicator element.
	ttyprogress.RunWith(p, bar, func(bar ttyprogress.Bar) {
		for i := 0; i <= 20; i++ {
			bar.Set(i * 5 * 5)
			time.Sleep(time.Millisecond * 500)
		}
	})

	// close context to indicate, that no new indicators are created anymore
	p.Close()

	// wait until all indicators are finished
	p.Wait(context.Background())
}

// --- end usage ---
