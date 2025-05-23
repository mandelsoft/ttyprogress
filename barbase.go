package ttyprogress

import (
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/object"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

type barGroupNotifier struct {
	started bool
}

var _ specs.GroupNotifier = (*barGroupNotifier)(nil)

func (n *barGroupNotifier) Add(b ProgressElement, p any) {
	eff := b.(*IntBarBase[IntBarImpl])
	defer eff.elem.Lock()()

	if !n.started {
		eff.elem.Protected().SetTotal(0)
		n.started = true
	}
	eff.elem.Protected().SetTotal(eff.elem.Protected().Total() + 1)
}

func (*barGroupNotifier) Done(b ProgressElement, p any) {
	eff := b.(*IntBarBase[IntBarImpl])
	eff.Incr()
}

////////////////////////////////////////////////////////////////////////////////

type IntBarInterface interface {
	ppi.BarInterface[int]
	Current() int
	CompletedPercent() float64
	Set(n int) bool
	Incr() bool
}

type IntBarImpl interface {
	ppi.BarImpl[int]
	IntBarInterface
}

type IntBarBase[T IntBarImpl] struct {
	*ppi.BarBase[T, int]
	elem *IntBarBaseImpl[T]
}

func (b *IntBarBase[T]) CompletedPercent() float64 {
	defer b.elem.Lock()()

	return b.elem.Protected().CompletedPercent()
}

func (b *IntBarBase[T]) Current() int {
	defer b.elem.Lock()()

	return b.elem.Protected().Current()
}

func (b *IntBarBase[T]) Set(n int) bool {
	defer b.elem.Lock()()

	return b.elem.Protected().Set(n)
}

func (b *IntBarBase[T]) Incr() bool {
	defer b.elem.Lock()()

	return b.elem.Protected().Incr()
}

type IntBarBaseImpl[T IntBarImpl] struct {
	*ppi.BarBaseImpl[T, int]

	current int
}

func newIntBar[T IntBarImpl](p Container, c specs.BarBaseConfiguration, total int, self object.Self[T, any]) (*IntBarBase[T], *IntBarBaseImpl[T], error) {
	e := &IntBarBaseImpl[T]{}
	o := &IntBarBase[T]{elem: e}

	if self == nil {
		// T must be IntBarImpl
		self = object.NewSelf[T, any](generics.Cast[T](e), o)
	}

	b, s, err := ppi.NewBarBase[T, int](self, p, c, total, nil)
	if err != nil {
		return nil, nil, err
	}
	e.BarBaseImpl = s
	o.BarBase = b
	return o, e, nil
}

// Set the current count of the bar. It returns ErrMaxCurrentReached when trying n exceeds the total value. This is atomic operation and concurrency safe.
func (b *IntBarBaseImpl[T]) Set(n int) bool {
	b.Start()

	if b.current >= b.Protected().Total() {
		return false
	}
	if n >= b.Protected().Total() {
		n = b.Protected().Total()
	}
	b.current = n
	b.Protected().Flush()
	return true
}

// Incr increments the current value by 1, time elapsed to current time and returns true. It returns false if the cursor has reached or exceeds total value.
func (b *IntBarBaseImpl[T]) Incr() bool {
	b.Protected().Start()

	if b.current == b.Protected().Total() {
		return false
	}

	n := b.current + 1
	b.current = n
	b.Protected().Flush()
	return true
}

func (b *IntBarBaseImpl[T]) IsFinished() bool {
	return b.current == b.Protected().Total()
}

// Current returns the current progress of the bar
func (b *IntBarBaseImpl[T]) Current() int {
	return b.current
}

func runeBytes(r rune) []byte {
	return []byte(string(r))
}

// CompletedPercent return the percent completed
func (b *IntBarBaseImpl[T]) CompletedPercent() float64 {
	return (float64(b.Current()) / float64(b.Total())) * 100.00
}
