package ppi

import (
	"bytes"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/specs"
)

type BarInterface[V any] interface {
	ProgressInterface
	CompletedPercent() float64
	Total() V
}

type BarImpl[V any] interface {
	ProgressImpl
	BarInterface[V]
	SetTotal(v V)
}

type BarBase[T BarImpl[V], V any] struct {
	*ProgressBase[T]
	elem *BarBaseImpl[T, V]
}

func (b *BarBase[T, V]) Total() V {
	defer b.elem.Lock()()

	return b.elem.Protected().Total()
}

type BarBaseImpl[T BarImpl[V], V any] struct {
	*ProgressBaseImpl[T]

	// total of the total  for the progress bar.
	total V

	// pending is the message shown before started
	pending string

	config specs.BarConfig

	// width is the width of the progress bar.
	width uint
}

func (b *BarBaseImpl[T, V]) Total() V {
	return b.total
}

func (b *BarBaseImpl[T, V]) SetTotal(v V) {
	b.total = v
}

func NewBarBase[T BarImpl[V], V any](self Self[T, any], p Container, c specs.BarBaseConfiguration, total V, closer func(), tick ...bool) (*BarBase[T, V], *BarBaseImpl[T, V], error) {
	e := &BarBaseImpl[T, V]{
		total:   total,
		width:   c.GetWidth(),
		config:  c.GetConfig(),
		pending: c.GetPending(),
	}

	b, s, err := NewProgressBase[T](self, p, c, 1, closer, general.Optional(tick...))
	if err != nil {
		return nil, nil, err
	}
	e.ProgressBaseImpl = s
	return &BarBase[T, V]{b, e}, e, nil
}

func (b *BarBaseImpl[T, V]) GetBarConfig() specs.BarConfig {
	return b.config
}

func (b *BarBaseImpl[T, V]) GetPending() string {
	return b.pending
}

func (b *BarBaseImpl[T, V]) GetWidth() uint {
	return b.width
}

func runeBytes(r rune) []byte {
	return []byte(string(r))
}

func (b *BarBaseImpl[T, V]) Visualize() (ttycolors.String, bool) {
	var buf bytes.Buffer

	if !b.IsStarted() {
		return specs.String(b.pending), false
	}
	// render visualization
	if b.width > 0 {
		if b.config.LeftEnd != ' ' {
			buf.Write(runeBytes(b.config.LeftEnd))
		}
		completedWidth := int(float64(b.width) * (b.Protected().CompletedPercent() / 100.00))
		// add fill and empty bits

		fill := string(b.config.Fill)
		_ = fill
		for i := 0; i < completedWidth; i++ {
			buf.Write(runeBytes(b.config.Fill))
		}
		if completedWidth > 0 {
			if completedWidth < int(b.width) {
				buf.Write(runeBytes(b.config.Head))
			}
		} else {
			buf.Write(runeBytes(b.config.Empty))
		}
		for i := 0; i < int(b.width)-completedWidth-1; i++ {
			buf.Write(runeBytes(b.config.Empty))
		}

		buf.Write(runeBytes(b.config.RightEnd))
	}
	return ttycolors.Sequence(buf.String()), b.Protected().IsFinished()
}
