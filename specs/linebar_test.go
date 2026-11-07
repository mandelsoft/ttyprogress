package specs_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"

	"github.com/mandelsoft/ttyprogress/specs"
)

// lineBar is a line bar definition for the tests.
type lineBar struct {
	specs.LineBarDefinition[*lineBar, int]
}

func newLineBar() *lineBar {
	d := &lineBar{}
	d.LineBarDefinition = specs.NewLineBarDefinition(specs.NewSelf(d), 100)
	return d
}

var _ = Describe("Line bar definition", func() {
	ctx := ttycolors.NewContext(true)
	render := func(f ttycolors.Format) string {
		return ctx.Sequence(f.String("x")).String()
	}

	It("marks the completed part reversed by default", func() {
		d := newLineBar()
		Expect(render(d.GetCompletedFormat())).To(Equal(ansi.Mode(ansi.MODE_REVERSE) + "x" + ansi.Mode(ansi.MODE_RESET_REVERSE)))
		Expect(d.GetRemainingFormat()).To(BeNil())
	})

	It("composes the format for the completed part", func() {
		d := newLineBar().SetCompletedFormat(ttycolors.FmtBgGreen, ttycolors.FmtUnderline)
		Expect(render(d.GetCompletedFormat())).To(Equal(ansi.Mode(ansi.COLOR_BG_GREEN, ansi.MODE_UNDERLINE) + "x" +
			ansi.Mode(ansi.COLOR_BG_DEFAULT, ansi.MODE_RESET_UNDERLINE)))
	})

	It("sets and resets the format for the remaining part", func() {
		d := newLineBar().SetRemainingFormat(ttycolors.FmtItalic)
		Expect(render(d.GetRemainingFormat())).To(Equal(ansi.Mode(ansi.MODE_ITALIC) + "x" + ansi.Mode(ansi.MODE_RESET_ITALIC)))
		Expect(d.SetRemainingFormat().GetRemainingFormat()).To(BeNil())
	})

	It("keeps the formats for duplicates", func() {
		d := newLineBar().SetCompletedFormat(ttycolors.FmtBold).SetRemainingFormat(ttycolors.FmtDim)
		dup := &lineBar{}
		dup.LineBarDefinition = d.LineBarDefinition.Dup(specs.NewSelf(dup))
		Expect(dup.GetCompletedFormat()).To(BeIdenticalTo(d.GetCompletedFormat()))
		Expect(dup.GetRemainingFormat()).To(BeIdenticalTo(d.GetRemainingFormat()))
	})
})
