package specs

import (
	"fmt"
	"strings"
	"time"

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

// Message provides a DecoratorFunc for a simple text message.
func Message(m string) DecoratorFunc {
	return func(element ElementState) any {
		return m
	}
}

// Variable provides a DecoratorFunc for a variable of a progress indicator.
func Variable(name string) DecoratorFunc {
	return func(element ElementState) any {
		return element.(VariableProvider).GetVariable(name)
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

func (d *formattedDecoratorDefinition) CreateDecorator(e ElementState) types.Decorator {
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
	textWindow
	rolling
	speed *Speed
}

var _ types.Ticker = (*scrollingText)(nil)

func (s *scrollingText) Tick() bool {
	if s.speed.Tick() {
		s.rolling.Incr()
		return true
	}
	return false
}

func (s *scrollingText) Decorate() any {
	return ttycolors.Sequence(s.window(s.offset))
}

type scrollingTextDef struct {
	text   string
	n      int // number of characters of text
	length int
	speed  int
}

func (s *scrollingTextDef) CreateDecorator(e ElementState) types.Decorator {
	if s.n <= s.length {
		return Message(s.text + strings.Repeat(" ", s.length-s.n)).CreateDecorator(e)
	}
	return &scrollingText{
		// the text followed by a gap
		textWindow: newTextWindow(s.text+" ", s.n+1, s.length),
		rolling:    rolling{n: s.n + 1},
		speed:      NewSpeed(s.speed),
	}
}

// ScrollingText provides a scrolling text decorator.
// Scrolling is only done, if the text is longer than the
// specified text size (in characters). Otherwise, static
// padded text is shown. Output formats of the text are kept,
// other escape sequences are removed.
func ScrollingText(text string, length int) DecoratorDefinition {
	t, n := formattedText(text)
	return &scrollingTextDef{
		length: length,
		text:   t,
		n:      n,
		speed:  3,
	}
}

////////////////////////////////////////////////////////////////////////////////

type phasesText struct {
	phases Phases
	speed  *Speed
}

var _ types.Ticker = (*phasesText)(nil)

func (s *phasesText) Tick() bool {
	if s.speed.Tick() {
		s.phases.Incr()
		return true
	}
	return false
}

func (s *phasesText) Decorate() any {
	return s.phases.Get()
}

type phasesDef struct {
	phases Phases
	speed  int
}

// PhasesDecoration provides a decorator bases on Phases.
// Every decorator created for an element uses its own copy
// of the phases.
func PhasesDecoration(phases Phases, speed int) DecoratorDefinition {
	return &phasesDef{
		phases: phases,
		speed:  speed,
	}
}

func (s *phasesDef) CreateDecorator(e ElementState) types.Decorator {
	return &phasesText{
		phases: s.phases.Dup(),
		speed:  NewSpeed(s.speed),
	}
}
