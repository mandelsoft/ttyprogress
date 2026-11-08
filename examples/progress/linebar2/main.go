package main

import (
	"os"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/units"
)

func main() {
	// --- begin bar ---
	p := ttyprogress.For(os.Stdout)
	// Underline colors require an extended terminal feature, which cannot
	// be detected and is therefore disabled by default. Without it, the
	// underline color of the completed part is ignored.
	p.GetTTYContext().EnableFeatures(ttycolors.FeatureUnderlineColor)

	bar, _ := ttyprogress.NewLineBar().
		SetTotal(500).
		PrependMessage("Downloading...").
		PrependElapsed().AppendCompleted().SetDecoratorFormat(ttycolors.FmtGreen).
		AppendFunc(ttyprogress.Amount(units.Bytes(1024))).
		SetCompletedFormat(ttycolors.FmtUnderline, ttycolors.FmtUnderlineIdColor(ansi.COLOR_ID_YELLOW)).
		Add(p)
	// --- end bar ---

	for i := 0; i <= 20; i++ {
		bar.Set(i * 5 * 5)
		time.Sleep(time.Millisecond * 500)
	}
}
