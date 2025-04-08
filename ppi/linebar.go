package ppi

import (
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/object"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
	"github.com/mandelsoft/ttyprogress/specs"
)

type LineBarBase[T BarImpl[V], V any] struct {
	*ProgressBase[T]
	elem *LineBarBaseImpl[T, V]
}

func (b *LineBarBase[T, V]) Total() V {
	defer b.elem.Lock()()

	return b.elem.Protected().Total()
}

func (b *LineBarBase[T, V]) SetTotal(v V) {
	defer b.elem.Lock()()

	b.elem.Protected().SetTotal(v)
}

////////////////////////////////////////////////////////////////////////////////

type LineBarBaseImpl[T BarImpl[V], V any] struct {
	*ProgressBaseImpl[T]

	resetOnFinished bool

	// total of the total  for the progress bar.
	total V

	// pending is the message shown before started
	pending string
}

func (b *LineBarBaseImpl[T, V]) Total() V {
	return b.total
}

func (b *LineBarBaseImpl[T, V]) SetTotal(v V) {
	b.total = v
}

func NewLineBarBase[T BarImpl[V], V any](self object.Self[T, any], p Container, c specs.LineBarConfiguration[V], closer func(), tick ...bool) (*LineBarBase[T, V], *LineBarBaseImpl[T, V], error) {
	e := &LineBarBaseImpl[T, V]{
		total:           c.GetTotal(),
		pending:         c.GetPending(),
		resetOnFinished: c.IsResetOnFinished(),
	}

	b, s, err := NewProgressBase[T](self, p, c, 1, closer, general.Optional(tick...))
	if err != nil {
		return nil, nil, err
	}
	e.ProgressBaseImpl = s
	return &LineBarBase[T, V]{b, e}, e, nil
}

func (b *LineBarBaseImpl[T, V]) GetPending() string {
	return b.pending
}

func (b *LineBarBaseImpl[T, V]) Visualize() (ttycolors.String, bool) {
	if !b.IsStarted() {
		return specs.String(b.pending), false
	}

	return ttycolors.Sequence(""), b.Protected().IsFinished()
}

func (b *LineBarBaseImpl[T, V]) Line() (string, bool) {
	l, done := b.ProgressBaseImpl.Line()

	if done && b.resetOnFinished {
		return l, done
	}

	length := ansi.CharLen(l)
	completedWidth := int(float64(length) * (b.Protected().CompletedPercent() / 100.00))

	first, second := ansi.SplitAt(l, completedWidth)
	return ttycolors.Sequence(ttycolors.Reverse(first), second).String(), done

}
