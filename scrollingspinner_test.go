package ttyprogress_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
	"github.com/mandelsoft/ttyprogress"
)

var _ = Describe("Scrolling spinner", func() {
	texts := func(text string, length, n int) []string {
		p := ttyprogress.NewScrollingSpinner(text, length).GetPhases()
		var r []string
		for i := 0; i < n; i++ {
			r = append(r, ansi.StripAll(p.Get().String()))
			p.Incr()
		}
		return r
	}

	It("scrolls text with multi byte characters", func() {
		Expect(texts("äöüß", 3, 6)).To(Equal([]string{"äöü", "öüß", "üß ", "ß ä", " äö", "äöü"}))
	})

	It("scrolls short text", func() {
		Expect(texts("äö", 3, 3)).To(Equal([]string{"äö ", "ö  ", "  ä"}))
	})

	It("handles a window larger than the text", func() {
		Expect(texts("ä", 8, 2)).To(HaveLen(2))
	})

	It("removes escape sequences", func() {
		in := ttycolors.NewContext(true).Bold("äöüß").String()
		Expect(texts(in, 3, 2)).To(Equal([]string{"äöü", "öüß"}))
	})
})
