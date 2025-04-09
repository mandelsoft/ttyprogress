package ttyprogress

import (
	"fmt"
	"io"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/blocks"
	"github.com/mandelsoft/ttyprogress/types"
	"github.com/mandelsoft/ttyprogress/units"
)

////////////////////////////////////////////////////////////////////////////////

// Message is a decorator for progress elements
// providing a static message.
func Message(m ...any) DecoratorFunc {
	return func(element ElementState) any {
		return ttycolors.Sequence(m...)
	}
}

// Amount is a decorator for Bar elements
// providing information about the current and total amount
// for the progress.
func Amount(unit ...units.Unit) DecoratorFunc {
	u := general.OptionalDefaulted(units.Plain, unit...)
	return func(e ElementState) any {
		if t, ok := e.(interface{ Total() int }); ok {
			return fmt.Sprintf("(%s/%s)", u(e.(Bar).Current()), u(t.Total()))
		}
		return fmt.Sprintf("(%s)", u(e.(Bar).Current()))
	}
}

// Processed is a decorator for Bar objects
// providing information about the current progress value (int).
func Processed(unit ...units.Unit) DecoratorFunc {
	u := general.OptionalDefaulted(units.Plain, unit...)
	return func(e ElementState) any {
		return fmt.Sprintf("(%s)", u(e.(interface{ Current() int }).Current()))
	}
}

// PercentTerminalSize return a width relative to to the terminal size.
func PercentTerminalSize(p uint) uint {
	x, _ := blocks.GetTerminalSize()

	if x == 0 {
		return 10
	}
	s := (uint(x) * p) / 100
	if s < 10 {
		return 10
	}
	return s
}

// ReserveTerminalSize provide a reasonable width
// reserving an amount of characters for predefined fixed
// content.
func ReserveTerminalSize(r uint) uint {
	x, _ := blocks.GetTerminalSize()
	if x == 0 {
		return 10
	}
	s := x - int(r)
	if s < 10 {
		return 10
	}
	return uint(s)
}

// SimpleProgress creates and displays a single progress element according
// to the given specification.
func SimpleProgress[T Element](w io.Writer, e ElementDefinition[T]) T {
	p := For(w)
	b, _ := e.Add(p)
	p.Close()
	return b
}

func AddElement[T Element](p Container, definition types.ElementDefinition[T]) (T, error) {
	return definition.Add(p)
}

// RunWith asynchronously executes a function for progressing a
// progress indication given by its definition.
// If the function exists the indicator is closed
// The runner must accept an Element interface matching the
// created element type.
func RunWith[E Element, D ElementDefinition[E], A Element](c Container, def D, runner func(e A)) (E, error) {
	return RunWithS[E, D](c, def, func(e E) { runner(generics.Cast[A](e)) })
}

func RunWithS[E Element, D ElementDefinition[E]](c Container, def D, runner func(e E)) (E, error) {
	elem, err := def.Add(c)
	if err != nil {
		return elem, err
	}
	go func() {
		runner(elem)
		elem.Close()
	}()
	return elem, nil
}

// SpecializedRunner converts a generic runner to a type specific runner.
func SpecializedRunner[E Element](r func(Element)) func(E) {
	return func(e E) { r(e) }
}

// GeneralizedRunner converts a type specific runner to a general runner.
func GeneralizedRunner[E Element](r func(E)) func(Element) {
	return func(e Element) { r(generics.Cast[E](e)) }
}
