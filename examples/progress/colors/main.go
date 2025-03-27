package main

import (
	"fmt"
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
		PrependDecorator(specs.ScrollingText("https://go.dev/dl/go1.24.0.linux-amd64.tar.gz", 20)).
		PrependMessage("...").
		SetDecoratorFormat(ttycolors.FmtRGBColor(255, 255, 0)).
		PrependElapsed().
		SetDecoratorFormat(ttycolors.FmtBold).
		AppendCompleted().
		AppendFunc(ttyprogress.Amount(units.Bytes(1024))).
		Add(p)

	text, _ := ttyprogress.NewText(3).
		SetFollowUpGap(" > ").
		SetTitleLine("some output").
		SetTitleFormat(ttycolors.FmtRed).
		SetViewFormat(ttycolors.FmtBlue, ttycolors.FmtItalic).
		Add(p)

	spinner, _ := ttyprogress.NewSpinner().
		SetPhases(specs.NewFormatPhases("working",
			ttycolors.FmtRGBColor(255, 0, 0),
			ttycolors.FmtRGBColor(128, 128, 0),
			ttycolors.FmtRGBColor(0, 255, 0),
			ttycolors.FmtRGBColor(0, 128, 128),
			ttycolors.FmtRGBColor(0, 0, 255),
			ttycolors.FmtRGBColor(128, 0, 128),
		)).
		AppendElapsed().
		Add(p)

	p.Close()

	go func() {
		spinner.Start()
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 1000)
			text.Write([]byte(fmt.Sprintf("line %d\n", i)))
		}
		text.Close()
		spinner.Close()
	}()

	for i := 0; i <= 20; i++ {
		bar.Set(i * 5 * 5)
		time.Sleep(time.Millisecond * 500)
	}
	bar.Close()

	p.Wait(nil)
}
