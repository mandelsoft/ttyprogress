package main

import (
	"os"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/specs"
	"github.com/mandelsoft/ttyprogress/units"
)

func main() {
	p := ttyprogress.For(os.Stdout).EnableColors()

	bar, _ := ttyprogress.NewBar().
		SetPredefined(10).
		SetProgressColor(ttycolors.FmtBlue).
		SetTotal(500).
		SetColor(ttycolors.FmtCyan).
		SetWidth(ttyprogress.PercentTerminalSize(30)).
		PrependMessage("Downloading...").
		SetDecoratorFormat(ttycolors.FmtRGBColor(255, 0, 0)).
		PrependDecorator(specs.PhasesDecoration(specs.NewBouncingTextFmtPhases("https://go.dev/dl/go1.24.0.linux-amd64.tar.gz", 5, ttycolors.FmtStrikeThrough), 2)).
		PrependMessage("...").
		SetDecoratorFormat(ttycolors.FmtRGBColor(255, 255, 0)).
		PrependElapsed().
		SetDecoratorFormat(ttycolors.FmtBold).
		AppendCompleted().
		AppendFunc(ttyprogress.Amount(units.Bytes(1024))).
		Add(p)

	p.Close()

	for i := 0; i <= 20; i++ {
		bar.Set(i * 5 * 5)
		time.Sleep(time.Millisecond * 500)
	}
	bar.Close()

	p.Wait(nil)
}
