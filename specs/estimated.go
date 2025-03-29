package specs

import (
	"time"

	"github.com/mandelsoft/goutils/stringutils"
)

type EstimatedInterface interface {
	BarBaseInterface[time.Duration]

	Total() time.Duration
	TimeEstimated() time.Duration

	IsFinished() bool
	Set(n time.Duration) bool
}

type EstimatedDefinition[T any] struct {
	BarBaseDefinition[T]

	total time.Duration
}

// NewEstimatedDefinition can be used to create a nested definition
// for a derived bar definition.
func NewEstimatedDefinition[T any](self Self[T], total time.Duration) EstimatedDefinition[T] {
	d := EstimatedDefinition[T]{total: total}
	d.BarBaseDefinition = NewBarBaseDefinition(self)
	return d
}

func (d *EstimatedDefinition[T]) Dup(s Self[T]) EstimatedDefinition[T] {
	dup := *d
	dup.BarBaseDefinition = d.BarBaseDefinition.Dup(s)
	return dup
}

// PrependEstimated prepends the time elapsed to the beginning of the bar
func (d *EstimatedDefinition[T]) PrependEstimated(offset ...int) T {
	d.PrependFunc(func(e ElementState) any {
		return estimatedTime(e)
	}, offset...)
	return d.Self()
}

// AppendEstimated appends the time elapsed to the beginning of the bar
func (d *EstimatedDefinition[T]) AppendEstimated(offset ...int) T {
	d.AppendFunc(func(e ElementState) any {
		return estimatedTime(e)
	}, offset...)
	return d.Self()
}

func (d *EstimatedDefinition[T]) SetTotal(v time.Duration) T {
	d.total = v
	return d.Self()
}

func (d *EstimatedDefinition[T]) GetTotal() time.Duration {
	return d.total
}

////////////////////////////////////////////////////////////////////////////////

func estimatedTime(e ElementState) string {
	s := ""
	if e.IsStarted() {
		p := e.(EstimatedInterface)
		s = PrettyTime(p.Total() - p.TimeElapsed())
	}
	return stringutils.PadLeft(s, 5, ' ')
}

////////////////////////////////////////////////////////////////////////////////

type EstimatedSpecification[T any] interface {
	BarBaseSpecification[T]
	SetTotal(v time.Duration) T
}

type EstimatedConfiguration interface {
	BarBaseConfiguration
	GetTotal() time.Duration
}
