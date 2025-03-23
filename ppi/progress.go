package ppi

import (
	"slices"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/specs"
)

// ProgressInterface in the public interface of elements
// featuring a concrete progress information line.
type ProgressInterface = specs.ProgressInterface

////////////////////////////////////////////////////////////////////////////////

// ProgressProtected in the (protected) implementation interface for progress
// indicators.
type ProgressProtected[I any] interface {
	ElementProtected[I]
	Visualize() (string, bool)
}

func ProgressSelf[I ProgressInterface](impl ProgressProtected[I]) Self[I, ProgressProtected[I]] {
	return NewSelf[I, ProgressProtected[I]](impl)
}

// ProgressBase is a base implementation for elements providing
// a line for progress information.
type ProgressBase[T ProgressInterface] struct {
	ElemBase[T, ProgressProtected[T]]

	format         ttycolors.Format
	progressFormat ttycolors.Format
	appendFuncs    []DecoratorFunc
	prependFuncs   []DecoratorFunc
	tick           bool
}

func NewProgressBase[T ProgressInterface](self Self[T, ProgressProtected[T]], p Container, c specs.ProgressConfiguration, view int, closer func(), tick ...bool) (*ProgressBase[T], error) {
	e := &ProgressBase[T]{tick: general.Optional(tick...)}
	e.format = c.GetColor()
	e.progressFormat = c.GetProgressColor()
	e.prependFuncs = slices.Clone(c.GetPrependFuncs())
	e.appendFuncs = slices.Clone(c.GetAppendFuncs())

	b, err := NewElemBase[T, ProgressProtected[T]](self, p, c, view, closer)
	if err != nil {
		return nil, err
	}
	e.ElemBase = *b
	e.tick = general.OptionalDefaulted(c.GetTick(), tick...)
	return e, nil
}

func (b *ProgressBase[T]) Tick() bool {
	if b.tick && !b.closed && b.IsStarted() {
		b.self.Protected().Update()
		return true
	}
	return false
}

// AppendFunc runs the decorator function and renders the output on the right of the progress bar
func (b *ProgressBase[T]) AppendFunc(f DecoratorFunc, offset ...int) T {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	if len(offset) == 0 {
		b.appendFuncs = append(b.appendFuncs, f)
	} else {
		b.appendFuncs = slices.Insert(b.appendFuncs, offset[0], f)
	}
	return b.self.Self()
}

func (b *ProgressBase[T]) Line() (string, bool) {
	b.Lock.RLock()
	defer b.Lock.RUnlock()

	seq := make([]any, 0, 30)
	sep := false

	// render prepend functions to the left of the bar
	for _, f := range b.prependFuncs {
		if sep {
			seq = append(seq, " ")
		}
		seq = append(seq, f(b.self.Self()))
		sep = true
	}

	data, done := b.self.Protected().Visualize()
	// render main function
	if len(data) > 0 {
		if sep {
			seq = append(seq, " ")
		}
		if b.progressFormat != nil {
			seq = append(seq, b.progressFormat.String(string(data)))
		} else {
			seq = append(seq, string(data))
		}
		sep = true
	}

	// render append functions to the right of the bar
	for _, f := range b.appendFuncs {
		if sep {
			seq = append(seq, " ")
		}
		seq = append(seq, f(b.self.Self()))
		sep = true
	}

	if b.format != nil {
		return b.format.String(seq...).String(), done
	}
	return ttycolors.Sequence(seq...).String(), done
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
