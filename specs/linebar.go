package specs

import (
	"github.com/mandelsoft/goutils/optionutils"
	"github.com/mandelsoft/ttycolors"
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
	completedFormat ttycolors.Format
	remainingFormat ttycolors.Format
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
		completedFormat:    ttycolors.FmtReverse,
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

// SetCompletedFormat sets the output format for the part of the
// progress line corresponding to the completed progress. The given
// formats are composed (see ttycolors.New). The default is
// ttycolors.FmtReverse. The format applies to all output formats
// not set by the content of the line, for example by decorator formats.
func (d *LineBarDefinition[T, V]) SetCompletedFormat(fmts ...ttycolors.FormatProvider) T {
	d.completedFormat = ttycolors.New(fmts...)
	return d.Self()
}

// GetCompletedFormat provides the output format for the completed part.
func (d *LineBarDefinition[T, V]) GetCompletedFormat() ttycolors.Format {
	return d.completedFormat
}

// SetRemainingFormat sets the output format for the part of the
// progress line corresponding to the remaining progress. The given
// formats are composed (see ttycolors.New). Without formats, the
// remaining part is not formatted, which is the default.
func (d *LineBarDefinition[T, V]) SetRemainingFormat(fmts ...ttycolors.FormatProvider) T {
	if len(fmts) == 0 {
		d.remainingFormat = nil
	} else {
		d.remainingFormat = ttycolors.New(fmts...)
	}
	return d.Self()
}

// GetRemainingFormat provides the output format for the remaining
// part, or nil, if it is not formatted.
func (d *LineBarDefinition[T, V]) GetRemainingFormat() ttycolors.Format {
	return d.remainingFormat
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
	SetCompletedFormat(fmts ...ttycolors.FormatProvider) T
	SetRemainingFormat(fmts ...ttycolors.FormatProvider) T
}

type LineBarConfiguration[T any] interface {
	ProgressConfiguration
	IsResetOnFinished() bool
	GetPending() string
	GetTotal() T
	GetCompletedFormat() ttycolors.Format
	GetRemainingFormat() ttycolors.Format
}
