package specs_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
	"github.com/mandelsoft/ttyprogress/specs"
	"github.com/mandelsoft/ttyprogress/types"
)

var _ = Describe("Scrolling text", func() {
	decorate := func(d types.Decorator) string {
		return fmt.Sprint(d.Decorate())
	}

	It("scrolls text with multi byte characters", func() {
		d := specs.ScrollingText("äöü", 2).CreateDecorator(nil)
		t := d.(types.Ticker)
		Expect(decorate(d)).To(Equal("äö"))
		Expect(t.Tick()).To(BeTrue())
		Expect(decorate(d)).To(Equal("öü"))
	})

	It("pads short text by characters", func() {
		d := specs.ScrollingText("äö", 4).CreateDecorator(nil)
		Expect(decorate(d)).To(Equal("äö  "))
	})

	It("keeps output formats and removes other escape sequences", func() {
		colored := ttycolors.NewContext(true)
		render := func(d types.Decorator) string {
			return colored.Sequence(d.Decorate()).String()
		}
		b := ansi.Mode(ansi.MODE_BOLD)
		nb := ansi.Mode(ansi.MODE_RESET_BOLD)
		in := colored.Bold("äöü").String() + "\x1b[2K"
		d := specs.ScrollingText(in, 2).CreateDecorator(nil)
		Expect(render(d)).To(Equal(b + "äö" + nb))
		// advance two steps
		for steps := 0; steps < 2; {
			if d.(types.Ticker).Tick() {
				steps++
			}
		}
		Expect(render(d)).To(Equal(b + "ü" + nb + " "))
		Expect(render(specs.ScrollingText(in, 4).CreateDecorator(nil))).To(Equal(b + "äöü" + nb + " "))
	})
})
