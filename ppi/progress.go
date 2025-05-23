package ppi

import (
	"fmt"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/object"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
	"github.com/mandelsoft/ttyprogress/specs"
	"github.com/mandelsoft/ttyprogress/types"
)

// ProgressInterface in the public interface of elements
// featuring a concrete progress information line.
type ProgressInterface = specs.ProgressInterface

type ProgressImpl interface {
	ElementImpl
	specs.ProgressInterface
	Line() (string, bool)
	Tick() bool
	/* abstract protected */ Visualize() (ttycolors.String, bool)
	IsAutoClose() bool
}

// ProgressBase is a base implementation for elements providing
// a line for progress information.
type ProgressBase[T ProgressImpl] struct {
	*ElemBase[T]
	// synchronized
	elem *ProgressBaseImpl[T]
}

func (b *ProgressBase[T]) SetProgressColor(fmt ttycolors.FormatProvider) {
	b.elem.lock.Lock()
	defer b.elem.lock.Unlock()
	b.elem.SetProgressColor(fmt)
}

func (b *ProgressBase[T]) SetVariable(name string, value any) {
	b.elem.lock.Lock()
	defer b.elem.lock.Unlock()
	b.elem.Protected().SetVariable(name, value)
}

func (b *ProgressBase[T]) GetVariable(name string) any {
	b.elem.lock.RLock()
	defer b.elem.lock.RUnlock()
	return b.elem.Protected().GetVariable(name)
}

func (b *ProgressBase[T]) Tick() bool {
	b.elem.lock.RLock()
	defer b.elem.lock.RUnlock()
	return b.elem.Protected().Tick()
}

type ProgressBaseImpl[T ProgressImpl] struct {
	*ElemBaseImpl[T]

	format            ttycolors.Format
	progressFormat    ttycolors.Format
	appendDecorators  []types.Decorator
	prependDecorators []types.Decorator
	variables         map[string]any
	autoclose         bool
	minColumn         int

	tick    bool
	tickers []types.Ticker
}

var _ ElementImpl = (*ProgressBaseImpl[ProgressImpl])(nil)

func NewProgressBase[T ProgressImpl](self object.Self[T, any], p Container, c specs.ProgressConfiguration, view int, closer func(), tick ...bool) (*ProgressBase[T], *ProgressBaseImpl[T], error) {
	e := &ProgressBaseImpl[T]{
		tick:           general.Optional(tick...),
		variables:      make(map[string]any),
		autoclose:      c.IsAutoClose(),
		minColumn:      c.GetMinVisualizationColumn(),
		format:         c.GetColor(),
		progressFormat: c.GetProgressColor(),
	}

	for _, def := range c.GetPrependDecorators() {
		d := def.CreateDecorator(self.Protected())
		if t, ok := generics.UnwrapUntil[types.Ticker](d); ok {
			e.tickers = append(e.tickers, t)
		}
		e.prependDecorators = append(e.prependDecorators, d)
	}
	for _, def := range c.GetAppendDecorators() {
		d := def.CreateDecorator(self.Protected())
		if t, ok := generics.UnwrapUntil[types.Ticker](d); ok {
			e.tickers = append(e.tickers, t)
		}
		e.appendDecorators = append(e.appendDecorators, d)
	}
	if len(e.tickers) > 0 {
		e.tick = true
	}
	b, s, err := NewElemBase[T](self, p, c, view, closer)
	if err != nil {
		return nil, nil, err
	}
	e.ElemBaseImpl = s
	e.tick = general.OptionalDefaulted(c.GetTick(), tick...)
	return &ProgressBase[T]{b, e}, e, nil
}

func (b *ProgressBaseImpl[T]) SetProgressColor(fmt ttycolors.FormatProvider) {
	if fmt == nil {
		b.progressFormat = nil
	} else {
		b.progressFormat = fmt.Format()
	}
	b.Protected().Flush()
}

func (b *ProgressBaseImpl[T]) SetVariable(name string, value any) {
	b.variables[name] = value
	b.Protected().Flush()
}

func (b *ProgressBaseImpl[T]) GetVariable(name string) any {
	return b.variables[name]
}

func (b *ProgressBaseImpl[T]) IsAutoClose() bool {
	return b.autoclose
}

func (b *ProgressBaseImpl[T]) Tick() bool {
	if b.tick && !b.closed && b.IsStarted() {
		upd := false
		for _, t := range b.tickers {
			upd = t.Tick() || upd
		}
		return b.Protected().Update() || upd
	}
	return false
}

func (b *ProgressBaseImpl[T]) Line() (string, bool) {
	seq := make([]any, 0, 30)
	sep := false

	// render prepend functions to the left of the bar
	seq, sep = appendDecorators(seq, sep, b.prependDecorators)

	if b.minColumn > 0 {
		l := ansi.CharLen(b.block.GetGap() + b.String(seq...).String())
		if l < b.minColumn {
			seq = append(seq, fmt.Sprintf("%*s", b.minColumn-l, ""))
		}
	}

	data, done := b.Protected().Visualize()
	// render main function
	if data != nil {
		if sep {
			seq = append(seq, " ")
		}
		if b.progressFormat != nil {
			seq = append(seq, b.progressFormat.String(data))
		} else {
			seq = append(seq, data)
		}
		sep = true
	}

	// render append functions to the right of the bar
	seq, sep = appendDecorators(seq, sep, b.appendDecorators)

	if b.format != nil {
		return b.StringWith(b.format, seq...).String(), done
	}
	return b.String(seq...).String(), done
}

func appendDecorators(seq []any, sep bool, decorators []types.Decorator) ([]any, bool) {
	for _, f := range decorators {
		v := f.Decorate()
		if v != nil && v != "" {
			if sep {
				seq = append(seq, " ")
			}
			seq = append(seq, v)
			sep = true
		}
	}
	return seq, sep
}

func (b *ProgressBaseImpl[T]) Update() bool {
	line, done := b.Protected().Line()

	b.block.Reset()
	b.block.Write([]byte(line + "\n"))
	b.block.Flush()
	if done {
		if b.Protected().IsAutoClose() {
			b.Close()
		}
	}
	return true
}
