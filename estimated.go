package ttyprogress

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mandelsoft/ttycolors"
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

// _Estimated represents a progress bar
type _Estimated struct {
	barBase[Estimated, time.Duration]
}

type _estimatedProtected struct {
	*_Estimated
}

func (b *_estimatedProtected) Self() Estimated {
	return b._Estimated
}

func (b *_estimatedProtected) Update() bool {
	return b._update()
}

func (b *_estimatedProtected) Visualize() (ttycolors.String, bool) {
	return b._visualize()
}

// newEstimated returns a new progress bar
// based on expected execution time.
func newEstimated(p Container, c specs.EstimatedConfiguration) (Estimated, error) {
	e := &_Estimated{
		barBase: barBase[Estimated, time.Duration]{
			total:   c.GetTotal(),
			width:   c.GetWidth(),
			config:  c.GetConfig(),
			pending: c.GetPending(),
		},
	}
	self := ppi.ProgressSelf[Estimated](&_estimatedProtected{e})
	b, err := ppi.NewProgressBase[Estimated](self, p, c, 1, e.close, true)
	if err != nil {
		return nil, err
	}
	e.ProgressBase = b
	return e, nil
}

func (b *_Estimated) TimeEstimated() time.Duration {
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
func (b *_Estimated) Set(n time.Duration) bool {
	b.Start()

	elapsed := b.TimeElapsed()
	b.Lock.Lock()
	b.total = n
	if elapsed >= b.total {
		b.total += time.Second
	}
	b.Lock.Unlock()
	b.Flush()
	return true
}

func (b *_Estimated) close() {
	elapsed := b.TimeElapsed()

	b.Lock.Lock()
	defer b.Lock.Unlock()

	b.total = elapsed
}

func (b *_Estimated) IsFinished() bool {
	return b.IsClosed()
}

func (b *_Estimated) Current() time.Duration {
	return b.TimeElapsed()
}

// Total returns the expected goal.
func (b *_Estimated) Total() time.Duration {
	b.Lock.RLock()
	defer b.Lock.RUnlock()
	return b.total
}

func (b *_Estimated) _update() bool {
	return ppi.Update[Estimated](b.ProgressBase)
}

func (b *_Estimated) _visualize() (ttycolors.String, bool) {
	var buf bytes.Buffer

	if !b.IsStarted() {
		return specs.String(b.pending), false
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
	return specs.String(buf.String()), b.IsClosed()
}

// CompletedPercent return the percent completed
func (b *_Estimated) CompletedPercent() float64 {
	elapsed := b.TimeElapsed()
	total := b.total
	if total <= elapsed {
		total = elapsed
	}
	return (float64(elapsed) / float64(total)) * 100.00
}

// CompletedPercentString returns the formatted string representation of the completed percent
func (b *_Estimated) CompletedPercentString() string {
	return fmt.Sprintf("%3.f%%", b.CompletedPercent())
}
