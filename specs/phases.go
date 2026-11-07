package specs

import (
	"slices"
	"strings"
	"sync"

	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
)

////////////////////////////////////////////////////////////////////////////////

type Cycle[T any] struct {
	lock  sync.Mutex
	elems []T
	cnt   int
}

func NewCycle[T any](elems []T) *Cycle[T] {
	return &Cycle[T]{elems: slices.Clone(elems)}
}

func (s *Cycle[T]) Incr() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.cnt++
	if s.cnt >= len(s.elems) {
		s.cnt = 0
	}
}

func (s *Cycle[T]) Get() T {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.elems[s.cnt]
}

// Dup provides an independent copy of the cycle in its initial state.
func (s *Cycle[T]) Dup() *Cycle[T] {
	return &Cycle[T]{elems: s.elems}
}

// Phases describes the phases of an animated element.
// Every element using phases must use a dedicated instance
// (see Dup), because the phases keep the actual state.
type Phases interface {
	Incr()
	Get() ttycolors.String
	// Dup provides an independent copy in its initial state.
	Dup() Phases
}

// formats provides the formats of the given format providers.
// The formats are evaluated once, nil providers are ignored.
func formats(fmts []ttycolors.FormatProvider) []ttycolors.Format {
	var r []ttycolors.Format
	for _, f := range fmts {
		if f != nil {
			r = append(r, f.Format())
		}
	}
	return r
}

// cyclePhases are phases cycling through a list of formatted texts.
type cyclePhases struct {
	*Cycle[ttycolors.String]
}

var _ Phases = (*cyclePhases)(nil)

func (p *cyclePhases) Dup() Phases {
	return &cyclePhases{p.Cycle.Dup()}
}

// NewStaticPhases provides phases cycling through the given texts.
// Without texts, a single empty phase is provided.
func NewStaticPhases(phases ...string) Phases {
	if len(phases) == 0 {
		return newStaticPhase()
	}
	return &cyclePhases{NewCycle[ttycolors.String](sliceutils.Transform(phases, String))}
}

// NewStaticFormattedPhases provides phases cycling through the given
// formatted texts. Without texts, a single empty phase is provided.
func NewStaticFormattedPhases(phases ...ttycolors.String) Phases {
	if len(phases) == 0 {
		return newStaticPhase()
	}
	return &cyclePhases{NewCycle(phases)}
}

////////////////////////////////////////////////////////////////////////////////

// staticPhase is a single phase, which never changes.
// It is used, if the parameters of phases do not allow an animation.
type staticPhase struct {
	elems []any
}

var _ Phases = (*staticPhase)(nil)

// newStaticPhase provides a single phase showing the given elements.
func newStaticPhase(elems ...any) Phases {
	return &staticPhase{elems: elems}
}

func (p *staticPhase) Incr() {
}

func (p *staticPhase) Get() ttycolors.String {
	return ttycolors.Sequence(p.elems...)
}

// Dup provides the phase itself, it has no state.
func (p *staticPhase) Dup() Phases {
	return p
}

////////////////////////////////////////////////////////////////////////////////

type nestedPhases struct {
	phases Phases
	fmts   *Cycle[ttycolors.Format]
}

var _ Phases = (*nestedPhases)(nil)

// NewNestedFormatPhases creates phases based on a combination
// of bases phases formatted by a sequence of output formats.
// The number of phases for both parts may be different.
// Without formats (nil providers are ignored), the base phases are used
// unformatted.
func NewNestedFormatPhases(p Phases, fmts ...ttycolors.FormatProvider) Phases {
	f := formats(fmts)
	if len(f) == 0 {
		return p
	}
	return &nestedPhases{
		phases: p,
		fmts:   NewCycle(f),
	}
}

func (p *nestedPhases) Incr() {
	p.phases.Incr()
	p.fmts.Incr()
}

func (p *nestedPhases) Get() ttycolors.String {
	return p.fmts.Get().String(p.phases.Get())
}

func (p *nestedPhases) Dup() Phases {
	return &nestedPhases{phases: p.phases.Dup(), fmts: p.fmts.Dup()}
}

////////////////////////////////////////////////////////////////////////////////

type textFmtPhases struct {
	*Cycle[ttycolors.Format]
	text string
}

var _ Phases = (*textFmtPhases)(nil)

// NewFormatPhases provides phases for text with cycled formats.
// Without formats (nil providers are ignored), the text is shown
// unformatted.
func NewFormatPhases(text string, fmts ...ttycolors.FormatProvider) Phases {
	f := formats(fmts)
	if len(f) == 0 {
		return newStaticPhase(text)
	}
	return &textFmtPhases{
		text:  text,
		Cycle: NewCycle(f),
	}
}

func (p *textFmtPhases) Get() ttycolors.String {
	return p.Cycle.Get().String(p.text)
}

func (p *textFmtPhases) Dup() Phases {
	return &textFmtPhases{Cycle: p.Cycle.Dup(), text: p.text}
}

////////////////////////////////////////////////////////////////////////////////

// formattedText prepares a text, which may contain output format
// sequences, for cutting it into parts. Escape sequences other than
// output format sequences are removed. The output format sequences are
// normalized (see ttycolors.Normalize), so that resets ending formats in
// the text do not affect formats applied to parts of the text, and
// formats still set at its end are reset. Parts of the text are cut with
// ttycolors.SubString, so that they keep their output format.
// It provides the text together with its number of characters.
func formattedText(text string) (string, int) {
	t := ttycolors.Normalize(stripControls(text))
	n := ansi.CharLen(t)
	return ttycolors.SubString(t, 0, n), n
}

// stripControls removes all escape sequences from a text,
// which are no output format sequences.
func stripControls(text string) string {
	var b strings.Builder
	data := []byte(text)
	for len(data) > 0 {
		l, c := ansi.ScanSequence(data)
		switch {
		case l == 0:
			l = 1
			b.WriteByte(data[0])
		case c&ansi.S_CONTROL == 0:
			b.Write(data[:l])
		}
		data = data[l:]
	}
	return b.String()
}

// textWindow is the common part of phases showing a window
// of a given length (in characters) moving over a formatted text
// (see formattedText).
type textWindow struct {
	text   string
	n      int // number of characters of text
	length int
}

func newTextWindow(text string, n int, length int) textWindow {
	return textWindow{text: text, n: n, length: length}
}

// sub provides the characters from start to end-1 of the text.
func (w *textWindow) sub(start, end int) string {
	return ttycolors.SubString(w.text, start, end)
}

// window provides the window starting at the given position. A window
// exceeding the end of the text is continued at the beginning of the text.
func (w *textWindow) window(start int) string {
	end := start + w.length
	if end <= w.n {
		return w.sub(start, end)
	}
	return w.sub(start, w.n) + w.sub(0, end-w.n)
}

// formatted provides the text with the window starting at the given
// position formatted with the given format. A window exceeding the
// end of the text is continued at the beginning of the text.
// Formats of the text are kept, the format applies to all output
// formats not set by the text.
func (w *textWindow) formatted(f ttycolors.Format, start int) ttycolors.String {
	end := start + w.length
	if end <= w.n {
		return ttycolors.Sequence(w.sub(0, start), f.String(w.sub(start, end)), w.sub(end, w.n))
	}
	end -= w.n
	return ttycolors.Sequence(f.String(w.sub(0, end)), w.sub(end, start), f.String(w.sub(start, w.n)))
}

// rolling is a position running cyclically through n positions.
type rolling struct {
	offset int
	n      int
}

func (r *rolling) Incr() {
	r.offset = (r.offset + 1) % r.n
}

// bouncing is a position moving back and forth between 0 and limit.
type bouncing struct {
	offset int
	delta  int
	limit  int
}

// newBouncing provides a bouncing position starting at the given
// position and moving down (delta -1) or up (delta 1).
func newBouncing(start, delta, limit int) bouncing {
	return bouncing{offset: start, delta: delta, limit: limit}
}

func (b *bouncing) Incr() {
	if b.offset == 0 {
		b.delta = 1
	} else if b.offset >= b.limit {
		b.delta = -1
	}
	b.offset += b.delta
}

////////////////////////////////////////////////////////////////////////////////

type rollingTextPhases struct {
	textWindow
	rolling
}

var _ Phases = (*rollingTextPhases)(nil)

// NewRollingTextPhases provides phases for text rolling through
// a window of the given length (in characters). A text shorter than
// the window is padded with spaces, so that it rolls through the
// complete window. For an empty text, the window is shown empty.
// For a length less than 1, an empty phase is provided.
// Output formats of the text are kept, other escape sequences are removed.
func NewRollingTextPhases(text string, length int) Phases {
	text, n := formattedText(text)
	switch {
	case length < 1:
		return newStaticPhase()
	case n == 0:
		return newStaticPhase(strings.Repeat(" ", length))
	case n < length:
		text += strings.Repeat(" ", length-n)
		n = length
	}
	return &rollingTextPhases{
		textWindow: newTextWindow(text, n, length),
		rolling:    rolling{n: n},
	}
}

func (p *rollingTextPhases) Dup() Phases {
	return NewRollingTextPhases(p.text, p.length)
}

func (p *rollingTextPhases) Get() ttycolors.String {
	return ttycolors.Sequence(p.window(p.offset))
}

////////////////////////////////////////////////////////////////////////////////

type rollingTextFmtPhases struct {
	textWindow
	rolling
	fmt ttycolors.Format
}

var _ Phases = (*rollingTextFmtPhases)(nil)

// NewRollingTextFmtPhases provides phases for text with a formatted part
// of the given length (in characters) rolling through the text.
// If the length is not less than the length of the text, the
// complete text is shown formatted. For a length less than 1
// or a nil format, the text is shown unformatted.
// Output formats of the text are kept, the format applies to all
// output formats not set by the text. Other escape sequences are removed.
func NewRollingTextFmtPhases(text string, length int, fmt ttycolors.FormatProvider) Phases {
	text, n := formattedText(text)
	switch {
	case length < 1 || fmt == nil:
		return newStaticPhase(text)
	case length >= n:
		return newStaticPhase(fmt.Format().String(text))
	}
	return &rollingTextFmtPhases{
		textWindow: newTextWindow(text, n, length),
		rolling:    rolling{n: n},
		fmt:        fmt.Format(),
	}
}

func (p *rollingTextFmtPhases) Dup() Phases {
	return NewRollingTextFmtPhases(p.text, p.length, p.fmt)
}

func (p *rollingTextFmtPhases) Get() ttycolors.String {
	return p.formatted(p.fmt, p.offset)
}

///////////////////////////////////////////////////////////////////////////////

type bouncingTextPhases struct {
	textWindow
	bouncing
	orig string // text without padding
}

var _ Phases = (*bouncingTextPhases)(nil)

// NewBouncingTextPhases provides phases for text bouncing within
// a window of the given length (in characters).
// If the length is not greater than the length of the text, there
// is no room to bounce and the text is shown statically.
// Output formats of the text are kept, other escape sequences are removed.
func NewBouncingTextPhases(text string, length int) Phases {
	text, n := formattedText(text)
	if length <= n {
		return newStaticPhase(text)
	}
	limit := length - n
	pad := strings.Repeat(" ", limit)
	return &bouncingTextPhases{
		// the text padded on both sides, the window never exceeds its end.
		textWindow: newTextWindow(pad+text+pad, n+2*limit, length),
		bouncing:   newBouncing(limit, -1, limit),
		orig:       text,
	}
}

func (p *bouncingTextPhases) Dup() Phases {
	return NewBouncingTextPhases(p.orig, p.length)
}

func (p *bouncingTextPhases) Get() ttycolors.String {
	return ttycolors.Sequence(p.window(p.offset))
}

///////////////////////////////////////////////////////////////////////////////

type bouncingTextFmtPhases struct {
	textWindow
	bouncing
	fmt ttycolors.Format
}

var _ Phases = (*bouncingTextFmtPhases)(nil)

// NewBouncingTextFmtPhases provides phases for text with a formatted
// part of the given length (in characters) bouncing within the text.
// If the length is not less than the length of the text, the
// complete text is shown formatted. For a length less than 1
// or a nil format, the text is shown unformatted.
// Output formats of the text are kept, the format applies to all
// output formats not set by the text. Other escape sequences are removed.
func NewBouncingTextFmtPhases(text string, length int, fmt ttycolors.FormatProvider) Phases {
	text, n := formattedText(text)
	switch {
	case length < 1 || fmt == nil:
		return newStaticPhase(text)
	case length >= n:
		return newStaticPhase(fmt.Format().String(text))
	}
	return &bouncingTextFmtPhases{
		textWindow: newTextWindow(text, n, length),
		bouncing:   newBouncing(0, 1, n-length),
		fmt:        fmt.Format(),
	}
}

func (p *bouncingTextFmtPhases) Dup() Phases {
	return NewBouncingTextFmtPhases(p.text, p.length, p.fmt)
}

func (p *bouncingTextFmtPhases) Get() ttycolors.String {
	return p.formatted(p.fmt, p.offset)
}
