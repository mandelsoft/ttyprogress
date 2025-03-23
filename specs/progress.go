package specs

import (
	"slices"

	"github.com/mandelsoft/goutils/stringutils"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/types"
)

type ProgressInterface interface {
	ElementInterface
}

////////////////////////////////////////////////////////////////////////////////

type ProgressDefinition[T any] struct {
	ElementDefinition[T]

	format              ttycolors.Format
	progressFormat      ttycolors.Format
	nextdecoratorFormat ttycolors.Format
	appendFuncs         []DecoratorFunc
	prependFuncs        []DecoratorFunc
	tick                bool
}

var (
	_ ProgressSpecification[any] = (*ProgressDefinition[any])(nil)
	_ ProgressConfiguration      = (*ProgressDefinition[any])(nil)
	_ Prepender                  = (*ProgressDefinition[any])(nil)
	_ Appender                   = (*ProgressDefinition[any])(nil)
)

func NewProgressDefinition[T any](self Self[T]) ProgressDefinition[T] {
	return ProgressDefinition[T]{ElementDefinition: NewElementDefinition(self)}
}

func (d *ProgressDefinition[T]) Dup(s Self[T]) ProgressDefinition[T] {
	dup := *d
	dup.ElementDefinition = d.ElementDefinition.Dup(s)
	dup.appendFuncs = slices.Clone(dup.appendFuncs)
	dup.prependFuncs = slices.Clone(dup.prependFuncs)
	return dup
}

// GetTick returns whether a tick is required.
func (d *ProgressDefinition[T]) GetTick() bool {
	return d.tick
}

// SetDecoratorFormat sets the output format for the next decorator.
func (d *ProgressDefinition[T]) SetDecoratorFormat(col ttycolors.Format) T {
	d.nextdecoratorFormat = col
	return d.Self()
}

// SetColor sets the output format for the progress indicator line
func (d *ProgressDefinition[T]) SetColor(col ttycolors.Format) T {
	d.format = col
	return d.Self()
}

func (d *ProgressDefinition[T]) GetColor() ttycolors.Format {
	return d.format
}

// SetProgressColor sets the output format for the progress indicator.
func (d *ProgressDefinition[T]) SetProgressColor(col ttycolors.Format) T {
	d.progressFormat = col
	return d.Self()
}

func (d *ProgressDefinition[T]) GetProgressColor() ttycolors.Format {
	return d.progressFormat
}

func format(fmt *ttycolors.Format, f DecoratorFunc) DecoratorFunc {
	if *fmt == nil {
		return f
	}
	eff := *fmt
	*fmt = nil
	return func(e ElementInterface) any {
		return (eff).String(f(e))
	}
}

// AppendFunc2 runs the decorator function and renders the output on the right of the progress indicator
func (d *ProgressDefinition[T]) AppendFunc2(f DecoratorFunc, offset ...int) {
	f = format(&d.nextdecoratorFormat, f)
	if len(offset) == 0 {
		d.appendFuncs = append(d.appendFuncs, f)
	} else {
		d.appendFuncs = slices.Insert(d.appendFuncs, offset[0], f)
	}
}

// AppendFunc runs the decorator function and renders the output on the right of the progress indicator
func (d *ProgressDefinition[T]) AppendFunc(f DecoratorFunc, offset ...int) T {
	d.AppendFunc2(f, offset...)
	return d.Self()
}

func (d *ProgressDefinition[T]) GetAppendFuncs() []DecoratorFunc {
	return slices.Clone(d.appendFuncs)
}

// PrependFunc2 runs decorator function and render the output left the progress indicator
func (d *ProgressDefinition[T]) PrependFunc2(f DecoratorFunc, offset ...int) {
	f = format(&d.nextdecoratorFormat, f)
	if len(offset) == 0 {
		d.prependFuncs = append(d.prependFuncs, f)
	} else {
		d.prependFuncs = slices.Insert(d.prependFuncs, offset[0], f)
	}
}

// PrependFunc runs decorator function and render the output left the progress indicator
func (d *ProgressDefinition[T]) PrependFunc(f DecoratorFunc, offset ...int) T {
	d.PrependFunc2(f, offset...)
	return d.Self()
}

func (d *ProgressDefinition[T]) GetPrependFuncs() []DecoratorFunc {
	return slices.Clone(d.prependFuncs)
}

// AppendElapsed appends the time elapsed to the progress indicator
func (d *ProgressDefinition[T]) AppendElapsed(offset ...int) T {
	d.tick = true
	return d.AppendFunc(func(e ElementInterface) any {
		return stringutils.PadLeft(e.TimeElapsedString(), 5, ' ')
	}, offset...)
}

// PrependElapsed prepends the time elapsed to the beginning of the indicator
func (d *ProgressDefinition[T]) PrependElapsed(offset ...int) T {
	d.tick = true
	return d.PrependFunc(func(e ElementInterface) any {
		return stringutils.PadLeft(e.TimeElapsedString(), 5, ' ')
	}, offset...)
}

// AppendMessage appends text to the progress indicator
func (d *ProgressDefinition[T]) AppendMessage(m string, offset ...int) T {
	return d.AppendFunc(Message(m), offset...)
}

// PrependMessage prepends text to the beginning of the indicator
func (d *ProgressDefinition[T]) PrependMessage(m string, offset ...int) T {
	return d.PrependFunc(Message(m), offset...)
}

////////////////////////////////////////////////////////////////////////////////

// ProgressSpecification is the configuration interface for progress indicators.
type ProgressSpecification[T any] interface {
	ElementSpecification[T]

	// SetColor set the color used for the progress line.
	SetColor(col ttycolors.Format) T

	// SetProgressColor set the color used for the progress visualization.
	SetProgressColor(col ttycolors.Format) T

	// SetDecoratorFormat set the output format for the next decorator.
	SetDecoratorFormat(col ttycolors.Format) T

	// AppendFunc adds a function providing some text appended
	// to the basic progress indicator.
	// If there are implicit settings, the offset can be used to
	// specify the index in the list of functions.
	AppendFunc(f DecoratorFunc, offset ...int) T

	// PrependFunc adds a function providing some text prepended
	// to the basic progress indicator.
	// If there are implicit settings, the offset can be used to
	// specify the index in the list of functions.
	PrependFunc(f DecoratorFunc, offset ...int) T

	// AppendElapsed appends the elapsed time of the action
	// or the duration of the action if the element is already closed.
	AppendElapsed(offset ...int) T

	// PrependElapsed appends the elapsed time of the action
	// or the duration of the action if the element is already closed.
	PrependElapsed(offset ...int) T

	// AppendMessage appends text to the progress indicator.
	AppendMessage(m string, offset ...int) T

	// PrependMessage prepends text to the beginning of the indicator
	PrependMessage(m string, offset ...int) T
}

// ProgressConfiguration provides access to the generic Progress configuration
type ProgressConfiguration interface {
	ElementConfiguration
	GetTick() bool

	GetColor() ttycolors.Format
	GetProgressColor() ttycolors.Format
	GetPrependFuncs() []DecoratorFunc
	GetAppendFuncs() []DecoratorFunc
}

////////////////////////////////////////////////////////////////////////////////

func TransferProgressConfig[D ProgressSpecification[T], T any](d D, c ProgressConfiguration) D {
	for _, e := range c.GetPrependFuncs() {
		d.PrependFunc(e)
	}
	for _, e := range c.GetAppendFuncs() {
		d.AppendFunc(e)
	}
	d.SetColor(c.GetColor())
	d.SetFinal(c.GetFinal())
	return d
}

type Prepender interface {
	PrependFunc2(f DecoratorFunc, offset ...int)
}

type Appender interface {
	AppendFunc2(f DecoratorFunc, offset ...int)
}

func AppendFunc[T ElementInterface](d types.ElementDefinition[T], f DecoratorFunc, offset ...int) bool {
	if a, ok := types.Unwrap(d).(Appender); ok {
		a.AppendFunc2(f, offset...)
		return true
	}
	return false
}

func PrependFunc[T ElementInterface](d types.ElementDefinition[T], f DecoratorFunc, offset ...int) bool {
	if a, ok := types.Unwrap(d).(Prepender); ok {
		a.PrependFunc2(f, offset...)
		return true
	}
	return false
}
