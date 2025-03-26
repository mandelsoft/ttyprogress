package specs

import (
	"fmt"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/types"
	"github.com/mandelsoft/ttyprogress/units"
)

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
	offset int
	cnt    int
	speed  int
	text   string
}

var _ types.Ticker

func (s *scrollingText) Tick() bool {
	if len(s.text) > s.length {
		s.cnt++
		if s.cnt >= s.speed {
			s.cnt = 0
			s.offset = (s.offset + 1) % (len(s.text) + s.length)
			return true
		}
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

func (s *scrollingTextDef) CreateDecorator(_ ElementInterface) types.Decorator {
	t := s.text
	if len(s.text) > s.length {
		t += " "
	}
	return &scrollingText{
		gap:    fmt.Sprintf("%-*s", s.length, " "),
		length: s.length,
		text:   t,
		speed:  s.speed,
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
