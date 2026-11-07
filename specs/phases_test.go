package specs_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
	"github.com/mandelsoft/ttyprogress/specs"
	"github.com/mandelsoft/ttyprogress/types"
)

var colored = ttycolors.NewContext(true)

// render renders the actual phase with enabled output formatting.
func render(p specs.Phases) string {
	return colored.Sequence(p.Get()).String()
}

// phases renders n phases.
func phases(p specs.Phases, n int) []string {
	var r []string
	for i := 0; i < n; i++ {
		r = append(r, render(p))
		p.Incr()
	}
	return r
}

// texts renders n phases without escape sequences.
func texts(p specs.Phases, n int) []string {
	var r []string
	for _, s := range phases(p, n) {
		r = append(r, ansi.StripAll(s))
	}
	return r
}

// bold provides the rendered string for a bold text.
func bold(s string) string {
	return ansi.Mode(ansi.MODE_RESET_BOLD, ansi.MODE_BOLD) + s + ansi.Mode(ansi.MODE_RESET_BOLD)
}

// italic provides the rendered string for an italic text.
func italic(s string) string {
	return ansi.Mode(ansi.MODE_ITALIC) + s + ansi.Mode(ansi.MODE_RESET_ITALIC)
}

// static checks that all phases render the expected string.
func static(p specs.Phases, expected string) {
	for _, s := range phases(p, 5) {
		Expect(s).To(Equal(expected))
	}
}

// provider is a format provider, which is no ttycolors.Format.
type provider struct {
	fmt ttycolors.FormatProvider
}

var _ ttycolors.FormatProvider = (*provider)(nil)

func (p *provider) Format() ttycolors.FormatInfo {
	return p.fmt.Format()
}

var _ = Describe("Phases", func() {
	Context("static phases", func() {
		It("cycles texts", func() {
			Expect(texts(specs.NewStaticPhases("a", "b"), 3)).To(Equal([]string{"a", "b", "a"}))
		})

		It("provides an empty phase without texts", func() {
			static(specs.NewStaticPhases(), "")
			static(specs.NewStaticFormattedPhases(), "")
		})
	})

	Context("format phases", func() {
		It("cycles formats", func() {
			p := specs.NewFormatPhases("x", ttycolors.FmtBold, ttycolors.FmtRed)
			Expect(phases(p, 3)).To(Equal([]string{
				bold("x"),
				ansi.Mode(ansi.COLOR_FG_RED) + "x" + ansi.Mode(ansi.COLOR_FG_DEFAULT),
				bold("x"),
			}))
		})

		It("shows the text unformatted without formats", func() {
			static(specs.NewFormatPhases("x"), "x")
		})

		It("uses the base phases without formats", func() {
			base := specs.NewStaticPhases("a")
			Expect(specs.NewNestedFormatPhases(base)).To(BeIdenticalTo(base))
		})
	})

	Context("rolling text", func() {
		It("rolls text through a window", func() {
			Expect(texts(specs.NewRollingTextPhases("abcde", 3), 6)).To(Equal([]string{"abc", "bcd", "cde", "dea", "eab", "abc"}))
		})

		It("pads text shorter than the window", func() {
			Expect(texts(specs.NewRollingTextPhases("ab", 4), 5)).To(Equal([]string{"ab  ", "b  a", "  ab", " ab ", "ab  "}))
		})

		It("shows an empty window for an empty text", func() {
			static(specs.NewRollingTextPhases("", 3), "   ")
		})

		It("is empty for an invalid length", func() {
			static(specs.NewRollingTextPhases("abc", 0), "")
			static(specs.NewRollingTextPhases("abc", -1), "")
		})
	})

	Context("rolling formatted text", func() {
		It("rolls a formatted part through the text", func() {
			p := specs.NewRollingTextFmtPhases("abcd", 2, ttycolors.FmtBold)
			Expect(phases(p, 5)).To(Equal([]string{
				bold("ab") + "cd",
				"a" + bold("bc") + "d",
				"ab" + bold("cd"),
				bold("a") + "bc" + bold("d"),
				bold("ab") + "cd",
			}))
		})

		It("formats the complete text for a length not less than the text", func() {
			static(specs.NewRollingTextFmtPhases("abc", 3, ttycolors.FmtBold), bold("abc"))
			static(specs.NewRollingTextFmtPhases("abc", 5, ttycolors.FmtBold), bold("abc"))
		})

		It("shows the text unformatted for an invalid length", func() {
			static(specs.NewRollingTextFmtPhases("abc", 0, ttycolors.FmtBold), "abc")
		})
	})

	Context("bouncing text", func() {
		It("bounces text within a window", func() {
			Expect(texts(specs.NewBouncingTextPhases("ab", 4), 6)).To(Equal([]string{"ab  ", " ab ", "  ab", " ab ", "ab  ", " ab "}))
		})

		It("shows the text statically without room to bounce", func() {
			static(specs.NewBouncingTextPhases("abc", 3), "abc")
			static(specs.NewBouncingTextPhases("abc", 2), "abc")
			static(specs.NewBouncingTextPhases("abc", 0), "abc")
		})
	})

	Context("bouncing formatted text", func() {
		It("bounces a formatted part within the text", func() {
			p := specs.NewBouncingTextFmtPhases("abcd", 2, ttycolors.FmtBold)
			Expect(phases(p, 5)).To(Equal([]string{
				bold("ab") + "cd",
				"a" + bold("bc") + "d",
				"ab" + bold("cd"),
				"a" + bold("bc") + "d",
				bold("ab") + "cd",
			}))
		})

		It("formats the complete text for a length not less than the text", func() {
			static(specs.NewBouncingTextFmtPhases("abc", 3, ttycolors.FmtBold), bold("abc"))
			static(specs.NewBouncingTextFmtPhases("abc", 5, ttycolors.FmtBold), bold("abc"))
		})

		It("shows the text unformatted for an invalid length", func() {
			static(specs.NewBouncingTextFmtPhases("abc", 0, ttycolors.FmtBold), "abc")
		})
	})

	Context("Dup", func() {
		DescribeTable("provides an independent copy in its initial state",
			func(create func() specs.Phases) {
				initial := phases(create(), 4)

				p := create()
				p.Incr()
				p.Incr()
				before := render(p)

				d := p.Dup()
				Expect(phases(d, 4)).To(Equal(initial))
				Expect(render(p)).To(Equal(before))
			},
			Entry("static phases", func() specs.Phases { return specs.NewStaticPhases("a", "b", "c") }),
			Entry("formatted phases", func() specs.Phases {
				return specs.NewStaticFormattedPhases(ttycolors.Sequence("a"), ttycolors.Bold("b"))
			}),
			Entry("format phases", func() specs.Phases { return specs.NewFormatPhases("x", ttycolors.FmtBold, ttycolors.FmtRed) }),
			Entry("nested phases", func() specs.Phases {
				return specs.NewNestedFormatPhases(specs.NewStaticPhases("a", "b", "c"), ttycolors.FmtBold, ttycolors.FmtRed)
			}),
			Entry("rolling text", func() specs.Phases { return specs.NewRollingTextPhases("abcde", 3) }),
			Entry("padded rolling text", func() specs.Phases { return specs.NewRollingTextPhases("ab", 4) }),
			Entry("rolling formatted text", func() specs.Phases { return specs.NewRollingTextFmtPhases("abcd", 2, ttycolors.FmtBold) }),
			Entry("bouncing text", func() specs.Phases { return specs.NewBouncingTextPhases("ab", 4) }),
			Entry("bouncing formatted text", func() specs.Phases { return specs.NewBouncingTextFmtPhases("abcd", 2, ttycolors.FmtBold) }),
			Entry("static phase", func() specs.Phases { return specs.NewRollingTextPhases("", 3) }),
		)

		It("provides the static phase itself", func() {
			p := specs.NewFormatPhases("x")
			Expect(p.Dup()).To(BeIdenticalTo(p))
		})
	})

	Context("decoration", func() {
		It("uses an own copy of the phases for every decorator", func() {
			def := specs.PhasesDecoration(specs.NewStaticPhases("a", "b"), 1)
			d1 := def.CreateDecorator(nil)
			d2 := def.CreateDecorator(nil)

			Expect(d1.(types.Ticker).Tick()).To(BeTrue())
			Expect(ansi.StripAll(colored.Sequence(d1.Decorate()).String())).To(Equal("b"))
			Expect(ansi.StripAll(colored.Sequence(d2.Decorate()).String())).To(Equal("a"))
		})
	})

	Context("format providers", func() {
		bp := &provider{ttycolors.FmtBold}

		It("accepts format providers, which are no formats", func() {
			Expect(phases(specs.NewFormatPhases("x", bp), 1)).To(Equal([]string{bold("x")}))
			Expect(phases(specs.NewNestedFormatPhases(specs.NewStaticPhases("a"), bp), 1)).To(Equal([]string{bold("a")}))
			Expect(phases(specs.NewRollingTextFmtPhases("abcd", 2, bp), 1)).To(Equal([]string{bold("ab") + "cd"}))
			Expect(phases(specs.NewBouncingTextFmtPhases("abcd", 2, bp), 1)).To(Equal([]string{bold("ab") + "cd"}))
		})

		It("ignores nil providers", func() {
			Expect(phases(specs.NewFormatPhases("x", nil, ttycolors.FmtBold, nil), 2)).To(Equal([]string{bold("x"), bold("x")}))
			static(specs.NewFormatPhases("x", nil), "x")
			base := specs.NewStaticPhases("a")
			Expect(specs.NewNestedFormatPhases(base, nil)).To(BeIdenticalTo(base))
		})

		It("shows the text unformatted for a nil format", func() {
			static(specs.NewRollingTextFmtPhases("abcd", 2, nil), "abcd")
			static(specs.NewBouncingTextFmtPhases("abcd", 2, nil), "abcd")
		})
	})

	Context("multi byte characters", func() {
		It("rolls text", func() {
			Expect(texts(specs.NewRollingTextPhases("äöü", 2), 4)).To(Equal([]string{"äö", "öü", "üä", "äö"}))
			Expect(texts(specs.NewRollingTextPhases("äö", 3), 3)).To(Equal([]string{"äö ", "ö ä", " äö"}))
		})

		It("rolls a formatted part", func() {
			p := specs.NewRollingTextFmtPhases("äöüß", 2, ttycolors.FmtBold)
			Expect(phases(p, 4)).To(Equal([]string{
				bold("äö") + "üß",
				"ä" + bold("öü") + "ß",
				"äö" + bold("üß"),
				bold("ä") + "öü" + bold("ß"),
			}))
		})

		It("bounces text", func() {
			Expect(texts(specs.NewBouncingTextPhases("äö", 3), 3)).To(Equal([]string{"äö ", " äö", "äö "}))
		})

		It("bounces a formatted part", func() {
			p := specs.NewBouncingTextFmtPhases("äöü", 2, ttycolors.FmtBold)
			Expect(phases(p, 3)).To(Equal([]string{bold("äö") + "ü", "ä" + bold("öü"), bold("äö") + "ü"}))
		})

		It("counts characters for the edge cases", func() {
			static(specs.NewBouncingTextPhases("äö", 2), "äö")
			static(specs.NewRollingTextFmtPhases("äö", 2, ttycolors.FmtBold), bold("äö"))
			static(specs.NewBouncingTextFmtPhases("äö", 2, ttycolors.FmtBold), bold("äö"))
		})
	})

	Context("formatted input", func() {
		// bold text with an embedded control string (window title)
		in := colored.Bold("äb").String() + "\x1b]0;title\acd"
		b := ansi.Mode(ansi.MODE_BOLD)
		nb := ansi.Mode(ansi.MODE_RESET_BOLD)
		i := ansi.Mode(ansi.MODE_ITALIC)
		ni := ansi.Mode(ansi.MODE_RESET_ITALIC)

		It("keeps output formats for rolling text", func() {
			Expect(phases(specs.NewRollingTextPhases(in, 2), 4)).To(Equal([]string{
				b + "äb" + nb, b + "b" + nb + "c", "cd", "d" + b + "ä" + nb}))
			static(specs.NewRollingTextPhases(in, 0), "")
		})

		It("keeps output formats for rolling formatted text", func() {
			Expect(phases(specs.NewRollingTextFmtPhases(in, 2, ttycolors.FmtItalic), 3)).To(Equal([]string{
				i + b + "äb" + nb + ni + "cd",
				b + "ä" + nb + i + b + "b" + nb + "c" + ni + "d",
				b + "äb" + nb + i + "cd" + ni,
			}))
			static(specs.NewRollingTextFmtPhases(in, 4, ttycolors.FmtItalic), i+b+"äb"+nb+"cd"+ni)
		})

		It("keeps output formats for bouncing text", func() {
			Expect(phases(specs.NewBouncingTextPhases(in, 5), 2)).To(Equal([]string{b + "äb" + nb + "cd ", " " + b + "äb" + nb + "cd"}))
			static(specs.NewBouncingTextPhases(in, 4), b+"äb"+nb+"cd")
		})

		It("keeps output formats for bouncing formatted text", func() {
			Expect(phases(specs.NewBouncingTextFmtPhases(in, 2, ttycolors.FmtItalic), 1)).To(Equal([]string{i + b + "äb" + nb + ni + "cd"}))
			static(specs.NewBouncingTextFmtPhases(in, 0, ttycolors.FmtItalic), b+"äb"+nb+"cd")
		})

		It("keeps the format visible after resets in the text", func() {
			red := ansi.Mode(ansi.COLOR_FG_RED)
			nored := ansi.Mode(ansi.COLOR_FG_DEFAULT)
			bb := ansi.Mode(ansi.MODE_RESET_BOLD, ansi.MODE_BOLD)
			p := specs.NewBouncingTextFmtPhases(red+"red"+ansi.Mode(ansi.MODE_RESET_ALL)+" rest", 4, ttycolors.FmtBold)
			Expect(phases(p, 5)).To(Equal([]string{
				bb + red + "red" + nored + " " + nb + "rest",
				red + "r" + nored + bb + red + "ed" + nored + " r" + nb + "est",
				red + "re" + nored + bb + red + "d" + nored + " re" + nb + "st",
				red + "red" + nored + bb + " res" + nb + "t",
				red + "red" + nored + " " + bb + "rest" + nb,
			}))
		})

		It("keeps output formats of duplicates", func() {
			p := specs.NewRollingTextPhases(in, 2)
			p.Incr()
			Expect(render(p.Dup())).To(Equal(b + "äb" + nb))
		})
	})
})
