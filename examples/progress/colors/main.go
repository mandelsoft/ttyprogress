package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress"
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

	p.Close()

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 1000)
			text.Write([]byte(fmt.Sprintf("line %d\n", i)))
		}
		text.Close()
	}()

	for i := 0; i <= 20; i++ {
		bar.Set(i * 5 * 5)
		time.Sleep(time.Millisecond * 500)
	}
	bar.Close()

	p.Wait(nil)
}
