package ttyprogress

import (
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

type barGroupNotifier struct {
	started bool
}

var _ specs.GroupNotifier[Bar] = (*barGroupNotifier)(nil)

func (n *barGroupNotifier) Add(b Bar, p any) {
	eff := b.(*IntBarBase[IntBarImpl])
	eff.elem.Lock()
	defer eff.elem.Unlock()

	if !n.started {
		eff.elem.SetTotal(0)
		n.started = true
	}
	eff.elem.SetTotal(eff.elem.Total() + 1)
}

func (*barGroupNotifier) Done(b Bar, p any) {
	b.Incr()
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
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.CompletedPercent()
}

func (b *IntBarBase[T]) Current() int {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.Current()
}

func (b *IntBarBase[T]) Set(n int) bool {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.Set(n)
}

func (b *IntBarBase[T]) Incr() bool {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.Incr()
}

type IntBarBaseImpl[T IntBarImpl] struct {
	*ppi.BarBaseImpl[T, int]

	current int
}

func newIntBar[T IntBarImpl](p Container, c specs.BarBaseConfiguration, total int, self func(*IntBarBaseImpl[T]) T) (*IntBarBase[T], *IntBarBaseImpl[T], error) {
	e := &IntBarBaseImpl[T]{}

	b, s, err := ppi.NewBarBase[T, int](self(e), p, c, total, nil)
	if err != nil {
		return nil, nil, err
	}
	e.BarBaseImpl = s
	return &IntBarBase[T]{b, e}, e, nil
}

// Set the current count of the bar. It returns ErrMaxCurrentReached when trying n exceeds the total value. This is atomic operation and concurrency safe.
func (b *IntBarBaseImpl[T]) Set(n int) bool {
	b.Start()

	if b.current >= b.Total() {
		return false
	}
	if n >= b.Total() {
		n = b.Total()
	}
	b.current = n
	b.Flush()
	return true
}

// Incr increments the current value by 1, time elapsed to current time and returns true. It returns false if the cursor has reached or exceeds total value.
func (b *IntBarBaseImpl[T]) Incr() bool {
	b.Start()

	if b.current == b.Self().Total() {
		return false
	}

	n := b.current + 1
	b.current = n

	if b.current == b.Self().Total() {
		b.Close()
	}
	b.Flush()
	return true
}

func (b *IntBarBaseImpl[T]) IsFinished() bool {
	return b.current == b.Self().Total()
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
