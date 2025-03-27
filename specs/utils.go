package specs

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/mandelsoft/goutils/stringutils"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/types"
	"github.com/mandelsoft/ttyprogress/units"
)

func String(s string) ttycolors.String {
	if s == "" {
		return nil
	}
	return ttycolors.Sequence(s)
}

// Message provide a DecoratorFunc for a simple text message.
func Message(m string) DecoratorFunc {
	return func(element ElementInterface) any {
		return m
	}
}

// PercentString returns the formatted string representation of the percent value.
func PercentString(p float64) string {
	return fmt.Sprintf("%3.f%%", p)
}

func PrettyTime(t time.Duration) string {
	if t == 0 {
		return ""
	}
	return units.Seconds(int(t.Truncate(time.Second) / time.Second))
}

////////////////////////////////////////////////////////////////////////////////

type formattedDecoratorDefinition struct {
	format ttycolors.Format
	def    DecoratorDefinition
}

func (d *formattedDecoratorDefinition) CreateDecorator(e ElementInterface) types.Decorator {
	return &formattedDecorator{d.format, d.def.CreateDecorator(e)}
}

type formattedDecorator struct {
	format ttycolors.Format
	deco   types.Decorator
}

func (f *formattedDecorator) Decorate() any {
	return f.format.String((f.deco).Decorate())
}

func (f *formattedDecorator) Unwrap() any {
	return f.deco
}

func Formatted(def DecoratorDefinition, f ...ttycolors.FormatProvider) DecoratorDefinition {
	if len(f) == 0 {
		return def
	}
	return &formattedDecoratorDefinition{ttycolors.New(f...), def}
}

////////////////////////////////////////////////////////////////////////////////

type scrollingText struct {
	gap    string
	length int
	text   string
	offset int
	speed  *Speed
}

var _ types.Ticker

func (s *scrollingText) Tick() bool {
	if s.speed.Tick() {
		s.offset = (s.offset + 1) % (len(s.text) + s.length)
		return true
	}
	return false
}

func (s *scrollingText) Decorate() any {
	return (s.text + s.text)[s.offset : s.offset+s.length]
}

type scrollingTextDef struct {
	text   string
	length int
	speed  int
}

func (s *scrollingTextDef) CreateDecorator(e ElementInterface) types.Decorator {
	t := s.text
	if len(s.text) <= s.length {
		return Message(stringutils.PadRight(t, s.length, ' ')).CreateDecorator(e)
	}
	t += " "
	return &scrollingText{
		gap:    fmt.Sprintf("%-*s", s.length, " "),
		length: s.length,
		text:   t,
		speed:  NewSpeed(s.speed),
	}
}

// ScrollingText provides a scrolling text decorator.
// Scrolling is only done, if the text is longer than the
// specified text size. Otherwise, static padded text is shown.
func ScrollingText(text string, length int) DecoratorDefinition {
	return &scrollingTextDef{
		length: length,
		text:   text,
		speed:  3,
	}
}

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

type Phases interface {
	Incr()
	Get() ttycolors.String
}

var _ Phases = (*Cycle[ttycolors.String])(nil)

func NewStaticPhases(phases ...string) Phases {
	return NewCycle[ttycolors.String](sliceutils.Transform(phases, String))
}

func NewStaticFormattedPhases(phases ...ttycolors.String) Phases {
	return NewCycle(phases)
}

////////////////////////////////////////////////////////////////////////////////

// NewFormatPhases creates Phases using the same text with a sequence of formats.
func NewFormatPhases(s string, fmts ...ttycolors.FormatProvider) Phases {
	var phases []ttycolors.String

	for _, f := range fmts {
		phases = append(phases, f.Format().String(s))
	}
	return NewStaticFormattedPhases(phases...)
}

////////////////////////////////////////////////////////////////////////////////

type nestedPhases struct {
	phases Phases
	fmts   *Cycle[ttycolors.Format]
}

// NewNestedFormatPhases creates phases based on a combination
// of bases phases formatted by a sequence of output formats.
// The number of phases for both parts may be different.
func NewNestedFormatPhases(p Phases, fmts ...ttycolors.FormatProvider) Phases {
	return &nestedPhases{
		phases: p,
		fmts:   NewCycle[ttycolors.Format](sliceutils.Transform(fmts, func(f ttycolors.FormatProvider) ttycolors.Format { return f.Format() })),
	}
}

func (p *nestedPhases) Incr() {
	p.phases.Incr()
	p.fmts.Incr()
}

func (p *nestedPhases) Get() ttycolors.String {
	return p.fmts.Get().String(p.phases.Get())
}
