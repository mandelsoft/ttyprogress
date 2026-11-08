package main

import (
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	// --- begin bar ---
	total := 10 * time.Second
	bar, _ := ttyprogress.NewEstimated(total).
		SetWidth(ttyprogress.ReserveTerminalSize(40)).
		SetPredefined(10).
		PrependFunc(ttyprogress.Message("Downloading...")).
		PrependEstimated().
		AppendCompleted().
		AppendElapsed().
		Add(p)
	// --- end bar ---
	bar.Start()
	p.Close()

	for i := 0; i <= 19; i++ {
		time.Sleep(time.Millisecond * 500)
		// Adjust expected duration
		total = total + 100*time.Millisecond
		bar.Set(total)
	}
	time.Sleep(time.Second * 2)
	bar.Close()
}
