package specs

// Ginkgo and Gomega are not dot-imported, because package specs
// declares identifiers like Done and Pending.
import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
)

var _ = ginkgo.Describe("Phases building blocks", func() {
	ginkgo.It("rolls cyclically", func() {
		r := rolling{n: 3}
		var offsets []int
		for i := 0; i < 5; i++ {
			offsets = append(offsets, r.offset)
			r.Incr()
		}
		gomega.Expect(offsets).To(gomega.Equal([]int{0, 1, 2, 0, 1}))
	})

	ginkgo.It("bounces between 0 and the limit", func() {
		b := newBouncing(0, 1, 2)
		var offsets []int
		for i := 0; i < 7; i++ {
			offsets = append(offsets, b.offset)
			b.Incr()
		}
		gomega.Expect(offsets).To(gomega.Equal([]int{0, 1, 2, 1, 0, 1, 2}))

		b = newBouncing(2, -1, 2)
		offsets = nil
		for i := 0; i < 5; i++ {
			offsets = append(offsets, b.offset)
			b.Incr()
		}
		gomega.Expect(offsets).To(gomega.Equal([]int{2, 1, 0, 1, 2}))
	})

	ginkgo.Context("text window", func() {
		w := newTextWindow("äöüß", 4, 2)
		render := func(s ttycolors.String) string {
			return ttycolors.NewContext(true).Sequence(s).String()
		}
		bold := func(s string) string {
			return ansi.Mode(ansi.MODE_RESET_BOLD, ansi.MODE_BOLD) + s + ansi.Mode(ansi.MODE_RESET_BOLD)
		}

		ginkgo.It("formats a window within the text", func() {
			gomega.Expect(render(w.formatted(ttycolors.FmtBold, 1))).To(gomega.Equal("ä" + bold("öü") + "ß"))
			gomega.Expect(render(w.formatted(ttycolors.FmtBold, 2))).To(gomega.Equal("äö" + bold("üß")))
		})

		ginkgo.It("continues a window exceeding the text at the beginning", func() {
			gomega.Expect(render(w.formatted(ttycolors.FmtBold, 3))).To(gomega.Equal(bold("ä") + "öü" + bold("ß")))
		})
	})

	ginkgo.It("removes control sequences, but keeps output format sequences", func() {
		gomega.Expect(stripControls("\x1b[2Ka\x1b[1mb\x1b]0;title\ac\x1b")).To(gomega.Equal("a\x1b[1mbc\x1b"))
	})

	ginkgo.It("prepares formatted texts", func() {
		t, n := formattedText("\x1b[2K\x1b[31mred\x1b[0m \x1b[1mbold")
		gomega.Expect(n).To(gomega.Equal(8))
		gomega.Expect(t).To(gomega.Equal(ansi.Mode(ansi.COLOR_FG_RED) + "red" + ansi.Mode(ansi.COLOR_FG_DEFAULT) + " " +
			ansi.Mode(ansi.MODE_BOLD) + "bold" + ansi.Mode(ansi.MODE_RESET_BOLD)))
	})
})
