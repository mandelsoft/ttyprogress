package specs

import (
	"github.com/mandelsoft/goutils/optionutils"
)

type LineBarInterface interface {
	BarBaseInterface[int]
	SetTotal(m int)

	Set(n int) bool
	Incr() bool
}

type LineBarDefinition[T, V any] struct {
	ProgressDefinition[T]
	pending         string
	autoclose       bool
	resetOnFinished bool
	total           V
}

var (
	_ LineBarSpecification[any, int] = (*LineBarDefinition[any, int])(nil)
	_ LineBarConfiguration[int]      = (*LineBarDefinition[any, int])(nil)
)

// NewLineBarDefinition can be used to create a nested definition
// for a derived line bar definition.
func NewLineBarDefinition[T, V any](s Self[T], total V) LineBarDefinition[T, V] {
	return LineBarDefinition[T, V]{
		ProgressDefinition: NewProgressDefinition[T](s),
		total:              total,
	}
}

func (d *LineBarDefinition[T, V]) Dup(s Self[T]) LineBarDefinition[T, V] {
	dup := *d
	dup.ProgressDefinition = d.ProgressDefinition.Dup(s)
	return dup
}

// AppendCompleted appends the completion percent to the progress bar
func (d *LineBarDefinition[T, V]) AppendCompleted(offset ...int) T {
	d.AppendFunc(func(b ElementState) any {
		return PercentString(b.(CompletedPercent).CompletedPercent())
	}, offset...)
	return d.Self()
}

// PrependCompleted prepends the percent completed to the progress bar
func (d *LineBarDefinition[T, V]) PrependCompleted(offset ...int) T {
	d.PrependFunc(func(b ElementState) any {
		return PercentString(b.(CompletedPercent).CompletedPercent())
	}, offset...)
	return d.Self()
}

func (d *LineBarDefinition[T, V]) ResetOnFinished(b ...bool) T {
	d.resetOnFinished = optionutils.BoolOption(b...)
	return d.Self()
}

func (d *LineBarDefinition[T, V]) IsResetOnFinished() bool {
	return d.resetOnFinished
}

func (d *LineBarDefinition[T, V]) SetPending(v string) T {
	d.pending = v
	return d.Self()
}

func (d *LineBarDefinition[T, V]) GetPending() string {
	return d.pending
}

func (d *LineBarDefinition[T, V]) SetTotal(v V) T {
	d.total = v
	return d.Self()
}

func (d *LineBarDefinition[T, V]) GetTotal() V {
	return d.total
}

////////////////////////////////////////////////////////////////////////////////

type LineBarSpecification[T, V any] interface {
	ProgressSpecification[T]
	AppendCompleted(offset ...int) T
	PrependCompleted(offset ...int) T
	ResetOnFinished(b ...bool) T
	SetPending(m string) T
	SetTotal(v V) T
}

type LineBarConfiguration[T any] interface {
	ProgressConfiguration
	IsResetOnFinished() bool
	GetPending() string
	GetTotal() T
}
