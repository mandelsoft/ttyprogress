package ttyprogress

import (
	"time"

	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

type Estimated interface {
	specs.EstimatedInterface
}

type EstimatedDefinition struct {
	specs.EstimatedDefinition[*EstimatedDefinition]
}

func NewEstimated(total time.Duration) *EstimatedDefinition {
	d := &EstimatedDefinition{}
	d.EstimatedDefinition = specs.NewEstimatedDefinition[*EstimatedDefinition](specs.NewSelf(d), total)
	return d
}

func (d *EstimatedDefinition) Dup() *EstimatedDefinition {
	dup := &EstimatedDefinition{}
	dup.EstimatedDefinition = d.EstimatedDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *EstimatedDefinition) Add(c Container) (Estimated, error) {
	return newEstimated(c, d)
}

////////////////////////////////////////////////////////////////////////////////

type EstimatedInterface = specs.EstimatedInterface

// _Estimated represents a progress bar
type _Estimated struct {
	*ppi.BarBase[*_EstimatedImpl, time.Duration]
	elem *_EstimatedImpl
}

func (e *_Estimated) Set(d time.Duration) bool {
	defer e.elem.Lock()()

	return e.elem.Protected().Set(d)
}

func (e *_Estimated) Current() time.Duration {
	defer e.elem.Lock()()

	return e.elem.Protected().Current()
}

func (e *_Estimated) TimeEstimated() time.Duration {
	defer e.elem.Lock()()

	return e.elem.Protected().TimeEstimated()
}

func (e *_Estimated) CompletedPercent() float64 {
	defer e.elem.Lock()()

	return e.elem.Protected().CompletedPercent()
}

// _Estimated represents a progress bar
type _EstimatedImpl struct {
	*ppi.BarBaseImpl[*_EstimatedImpl, time.Duration]
}

// newEstimated returns a new progress bar
// based on expected execution time.
func newEstimated(p Container, c specs.EstimatedConfiguration) (Estimated, error) {
	e := &_EstimatedImpl{}
	o := &_Estimated{elem: e}

	b, s, err := ppi.NewBarBase[*_EstimatedImpl, time.Duration](ppi.NewSelf[*_EstimatedImpl, any](e, o), p, c, 1, e.closer, true)
	if err != nil {
		return nil, err
	}
	e.BarBaseImpl = s
	o.BarBase = b
	return o, nil
}

func (b *_EstimatedImpl) TimeEstimated() time.Duration {
	if b.IsStarted() {
		t := b.Total() - b.TimeElapsed()
		if t < 0 {
			if b.IsFinished() {
				t = 0
			} else {
				t = time.Second
			}
		}
		return t
	}
	return b.Total()
}

// Set the current count of the bar. It returns ErrMaxCurrentReached when trying n exceeds the total value. This is atomic operation and concurrency safe.
func (b *_EstimatedImpl) Set(n time.Duration) bool {
	b.Start()

	elapsed := b.TimeElapsed()
	b.SetTotal(n)
	if elapsed >= b.Total() {
		b.SetTotal(b.Total() + 2*time.Second)
	}
	b.Flush()
	return true
}

func (b *_EstimatedImpl) closer() {
	elapsed := b.TimeElapsed()
	b.SetTotal(elapsed)
}

func (b *_EstimatedImpl) IsFinished() bool {
	return b.IsClosed()
}

func (b *_EstimatedImpl) Current() time.Duration {
	return b.TimeElapsed()
}

// CompletedPercent return the percent completed
func (b *_EstimatedImpl) CompletedPercent() float64 {
	elapsed := b.TimeElapsed()
	total := b.Total()
	if total <= elapsed {
		total = elapsed
	}
	return (float64(elapsed) / float64(total)) * 100.00
}
