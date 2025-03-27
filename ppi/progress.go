package ppi

import (
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/specs"
	"github.com/mandelsoft/ttyprogress/types"
)

// ProgressInterface in the public interface of elements
// featuring a concrete progress information line.
type ProgressInterface = specs.ProgressInterface

////////////////////////////////////////////////////////////////////////////////

// ProgressProtected in the (protected) implementation interface for progress
// indicators.
type ProgressProtected[I any] interface {
	ElementProtected[I]
	Visualize() (ttycolors.String, bool)
}

func ProgressSelf[I ProgressInterface](impl ProgressProtected[I]) Self[I, ProgressProtected[I]] {
	return NewSelf[I, ProgressProtected[I]](impl)
}

// ProgressBase is a base implementation for elements providing
// a line for progress information.
type ProgressBase[T ProgressInterface] struct {
	*ElemBase[T, ProgressProtected[T]]

	format            ttycolors.Format
	progressFormat    ttycolors.Format
	appendDecorators  []types.Decorator
	prependDecorators []types.Decorator
	tick              bool
	tickers           []types.Ticker
}

func NewProgressBase[T ProgressInterface](self Self[T, ProgressProtected[T]], p Container, c specs.ProgressConfiguration, view int, closer func(), tick ...bool) (*ProgressBase[T], error) {
	e := &ProgressBase[T]{tick: general.Optional(tick...)}
	e.format = c.GetColor()
	e.progressFormat = c.GetProgressColor()

	for _, def := range c.GetPrependDecorators() {
		d := def.CreateDecorator(self.Self())
		if t, ok := generics.UnwrapUntil[types.Ticker](d); ok {
			e.tickers = append(e.tickers, t)
		}
		e.prependDecorators = append(e.prependDecorators, d)
	}
	for _, def := range c.GetAppendDecorators() {
		d := def.CreateDecorator(self.Self())
		if t, ok := generics.UnwrapUntil[types.Ticker](d); ok {
			e.tickers = append(e.tickers, t)
		}
		e.appendDecorators = append(e.appendDecorators, d)
	}
	if len(e.tickers) > 0 {
		e.tick = true
	}
	b, err := NewElemBase[T, ProgressProtected[T]](self, p, c, view, closer)
	if err != nil {
		return nil, err
	}
	e.ElemBase = b
	e.tick = general.OptionalDefaulted(c.GetTick(), tick...)
	return e, nil
}

func (b *ProgressBase[T]) Tick() bool {
	if b.tick && !b.closed && b.IsStarted() {
		upd := false
		for _, t := range b.tickers {
			upd = t.Tick() || upd
		}
		return b.self.Protected().Update() || upd
	}
	return false
}

func (b *ProgressBase[T]) Line() (string, bool) {
	b.Lock.RLock()
	defer b.Lock.RUnlock()

	seq := make([]any, 0, 30)
	sep := false

	// render prepend functions to the left of the bar
	for _, f := range b.prependDecorators {
		if sep {
			seq = append(seq, " ")
		}
		seq = append(seq, f.Decorate())
		sep = true
	}

	data, done := b.self.Protected().Visualize()
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
	for _, f := range b.appendDecorators {
		if sep {
			seq = append(seq, " ")
		}
		seq = append(seq, f.Decorate())
		sep = true
	}

	if b.format != nil {
		return b.StringWith(b.format, seq...).String(), done
	}
	return b.String(seq...).String(), done
}

func Update[T ProgressInterface](b *ProgressBase[T]) bool {
	line, done := b.Line()

	b.block.Reset()
	b.block.Write([]byte(line + "\n"))
	if done {
		b.Close()
	}
	return true
}

func (b *ProgressBase[T]) Flush() error {
	b.self.Protected().Update()
	return b.block.Flush()
}
