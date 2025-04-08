package ttyprogress

import (
	"github.com/mandelsoft/object"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

////////////////////////////////////////////////////////////////////////////////

// LineBar is a progress bar used to visualize the progress of an action in
// relation to a known maximum of required work reversing the progress line
// according to the actual progress.
type LineBar interface {
	specs.BarInterface
}

type LineBarDefinition struct {
	specs.LineBarDefinition[*LineBarDefinition, int]
}

var _ specs.GroupProgressElementDefinition[LineBar] = (*LineBarDefinition)(nil)

func NewLineBar() *LineBarDefinition {
	d := &LineBarDefinition{}
	d.LineBarDefinition = specs.NewLineBarDefinition(specs.NewSelf(d), 100)
	return d
}

func (d *LineBarDefinition) Dup() *LineBarDefinition {
	dup := &LineBarDefinition{}
	dup.LineBarDefinition = d.LineBarDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *LineBarDefinition) GetGroupNotifier() specs.GroupNotifier {
	return &barGroupNotifier{}
}

func (d *LineBarDefinition) Add(c Container) (LineBar, error) {
	return newLineBar(c, d)
}

////////////////////////////////////////////////////////////////////////////////

type _LineBar struct {
	*ppi.LineBarBase[*_LineBarImpl, int]
	elem *_LineBarImpl
}

func (b *_LineBar) CompletedPercent() float64 {
	defer b.elem.Lock()()

	return b.elem.Protected().CompletedPercent()
}

func (b *_LineBar) Current() int {
	defer b.elem.Lock()()

	return b.elem.Protected().Current()
}

func (b *_LineBar) SetTotal(n int) {
	defer b.elem.Lock()()

	b.elem.Protected().SetTotal(n)
}

func (b *_LineBar) Set(n int) bool {
	defer b.elem.Lock()()

	return b.elem.Protected().Set(n)
}

func (b *_LineBar) Incr() bool {
	defer b.elem.Lock()()

	return b.elem.Protected().Incr()
}

type _LineBarImpl struct {
	*ppi.LineBarBaseImpl[*_LineBarImpl, int]

	current int
}

func newLineBar(p Container, c specs.LineBarConfiguration[int]) (*_LineBar, error) {
	e := &_LineBarImpl{}
	o := &_LineBar{elem: e}

	b, s, err := ppi.NewLineBarBase[*_LineBarImpl, int](object.NewSelf[*_LineBarImpl, any](e, o), p, c, nil)
	if err != nil {
		return nil, err
	}
	e.LineBarBaseImpl = s
	o.LineBarBase = b
	return o, nil
}

// Set the current count of the bar. It returns ErrMaxCurrentReached when trying n exceeds the total value. This is atomic operation and concurrency safe.
func (b *_LineBarImpl) Set(n int) bool {
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
func (b *_LineBarImpl) Incr() bool {
	b.Protected().Start()

	if b.current == b.Protected().Total() {
		return false
	}

	n := b.current + 1
	b.current = n
	b.Protected().Flush()
	return true
}

func (b *_LineBarImpl) IsFinished() bool {
	return b.current == b.Protected().Total()
}

// Current returns the current progress of the bar
func (b *_LineBarImpl) Current() int {
	return b.current
}

// CompletedPercent return the percent completed
func (b *_LineBarImpl) CompletedPercent() float64 {
	return (float64(b.Current()) / float64(b.Total())) * 100.00
}
