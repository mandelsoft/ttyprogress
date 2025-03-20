package ttyprogress

import (
	"bytes"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

type BarConfig = specs.BarConfig
type Brackets = specs.Brackets

var (
	BarTypes     = specs.BarTypes
	BracketTypes = specs.BracketTypes

	// BarWidth is the default width of the progress bar
	BarWidth = specs.BarWidth

	// ErrMaxCurrentReached is error when trying to set current value that exceeds the total value
	ErrMaxCurrentReached = errors.New("errors: current value is greater total value")
)

// Bar is a progress bar used to visualize the progress of an action in
// relation to a known maximum of required work.
type Bar interface {
	specs.BarInterface
}

type BarDefinition struct {
	specs.BarDefinition[*BarDefinition]
}

var _ specs.GroupProgressElementDefinition[Bar] = (*BarDefinition)(nil)

func NewBar(set ...int) *BarDefinition {
	d := &BarDefinition{}
	d.BarDefinition = specs.NewBarDefinition(specs.NewSelf(d))
	if len(set) > 0 {
		d.SetPredefined(set[0])
	}
	return d
}

func (d *BarDefinition) Dup() *BarDefinition {
	dup := &BarDefinition{}
	dup.BarDefinition = d.BarDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *BarDefinition) GetGroupNotifier() specs.GroupNotifier[Bar] {
	return &barGroupNotifier{}
}

func (d *BarDefinition) Add(c Container) (Bar, error) {
	return newBar(c, d)
}

func (d *BarDefinition) AddWithTotal(c Container, total int) (Bar, error) {
	return newBar(c, d, total)
}

type barGroupNotifier struct {
	started bool
}

var _ specs.GroupNotifier[Bar] = (*barGroupNotifier)(nil)

func (n *barGroupNotifier) Add(b Bar, p any) {
	eff := b.(*_Bar[Bar])
	eff.Lock.Lock()
	defer eff.Lock.Unlock()
	if !n.started {
		eff.total = 0
		n.started = true
	}
	eff.total++
}

func (*barGroupNotifier) Done(b Bar, p any) {
	b.Incr()
}

////////////////////////////////////////////////////////////////////////////////

type barBase[T ppi.ProgressInterface, V any] struct {
	ppi.ProgressBase[T]

	// total of the total  for the progress bar.
	total V

	// pending is the message shown before started
	pending string

	config BarConfig

	// width is the width of the progress bar.
	width uint
}

func (b *barBase[T, V]) Total() V {
	b.Lock.RLock()
	defer b.Lock.RUnlock()

	return b.total
}

////////////////////////////////////////////////////////////////////////////////

// _Bar represents a progress bar
type _Bar[T ppi.ProgressInterface] struct {
	barBase[T, int]

	current int
}

type _barProtected struct {
	*_Bar[Bar]
}

func (b *_barProtected) Self() Bar {
	return b._Bar
}

func (b *_barProtected) Update() bool {
	return b._update()
}

func (b *_barProtected) Visualize() (string, bool) {
	return b._visualize()
}

// newBar returns a new progress bar
func newBar(p Container, c specs.BarConfiguration[int], total ...int) (Bar, error) {
	return newBarBase[Bar](p, c, general.OptionalDefaulted(c.GetTotal(), total...), func(e *_Bar[Bar]) ppi.Self[Bar, ppi.ProgressProtected[Bar]] {
		return ppi.ProgressSelf[Bar](&_barProtected{e})
	})
}

func newBarBase[T ppi.ProgressInterface](p Container, c specs.BarBaseConfiguration, total int, self func(*_Bar[T]) ppi.Self[T, ppi.ProgressProtected[T]]) (*_Bar[T], error) {
	e := &_Bar[T]{
		barBase: barBase[T, int]{
			total:   total,
			width:   c.GetWidth(),
			config:  c.GetConfig(),
			pending: c.GetPending(),
		},
	}
	b, err := ppi.NewProgressBase[T](self(e), p, c, 1, nil)
	if err != nil {
		return nil, err
	}
	e.ProgressBase = *b
	return e, nil
}

// Set the current count of the bar. It returns ErrMaxCurrentReached when trying n exceeds the total value. This is atomic operation and concurrency safe.
func (b *_Bar[T]) Set(n int) bool {
	b.Start()

	b.Lock.Lock()
	if b.current >= b.total {
		b.Lock.Unlock()
		return false
	}
	if n >= b.total {
		n = b.total
	}
	b.current = n
	b.Lock.Unlock()
	b.Flush()
	return true
}

// Incr increments the current value by 1, time elapsed to current time and returns true. It returns false if the cursor has reached or exceeds total value.
func (b *_Bar[T]) Incr() bool {
	b.Start()
	b.Lock.Lock()

	if b.incr() {
		if b.current == b.total {
			b.Lock.Unlock()
			b.Close()
		} else {
			b.Lock.Unlock()
			b.Flush()
		}
		return true
	}
	b.Lock.Unlock()
	return false
}

func (b *_Bar[T]) IsFinished() bool {
	b.Lock.RLock()
	defer b.Lock.RUnlock()
	return b.current == b.total
}

func (b *_Bar[T]) incr() bool {
	if b.current == b.total {
		return false
	}

	n := b.current + 1
	b.current = n
	return true
}

// Current returns the current progress of the bar
func (b *_Bar[T]) Current() int {
	b.Lock.RLock()
	defer b.Lock.RUnlock()
	return b.current
}

func (b *_Bar[T]) _update() bool {
	return ppi.Update[T](&b.ProgressBase)
}

func runeBytes(r rune) []byte {
	return []byte(string(r))
}

func (b *_Bar[T]) _visualize() (string, bool) {
	var buf bytes.Buffer

	if !b.IsStarted() {
		return b.pending, false
	}
	// render visualization
	if b.width > 0 {
		if b.config.LeftEnd != ' ' {
			buf.Write(runeBytes(b.config.LeftEnd))
		}
		completedWidth := int(float64(b.width) * (b.CompletedPercent() / 100.00))
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
	return buf.String(), b.current == b.total
}

// CompletedPercent return the percent completed
func (b *_Bar[T]) CompletedPercent() float64 {
	return (float64(b.Current()) / float64(b.total)) * 100.00
}
