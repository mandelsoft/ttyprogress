package specs

import (
	"slices"

	"github.com/mandelsoft/goutils/optionutils"
	"github.com/mandelsoft/goutils/stringutils"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/types"
)

type ProgressInterface = types.ProgressElement

////////////////////////////////////////////////////////////////////////////////

type ProgressDefinition[T any] struct {
	ElementDefinition[T]

	format              ttycolors.Format
	progressFormat      ttycolors.Format
	nextdecoratorFormat ttycolors.Format
	appendDefs          []DecoratorDefinition
	prependDefs         []DecoratorDefinition
	autoclose           bool
	minColumn           int
	tick                bool
}

var (
	_ ProgressSpecification[any] = (*ProgressDefinition[any])(nil)
	_ ProgressConfiguration      = (*ProgressDefinition[any])(nil)
	_ Prepender                  = (*ProgressDefinition[any])(nil)
	_ Appender                   = (*ProgressDefinition[any])(nil)
)

func NewProgressDefinition[T any](self Self[T]) ProgressDefinition[T] {
	return ProgressDefinition[T]{ElementDefinition: NewElementDefinition(self), autoclose: true}
}

func (d *ProgressDefinition[T]) Dup(s Self[T]) ProgressDefinition[T] {
	dup := *d
	dup.ElementDefinition = d.ElementDefinition.Dup(s)
	dup.appendDefs = slices.Clone(dup.appendDefs)
	dup.prependDefs = slices.Clone(dup.prependDefs)
	return dup
}

// GetTick returns whether a tick is required.
func (d *ProgressDefinition[T]) GetTick() bool {
	return d.tick
}

// setTick returns whether a tick is required.
func (d *ProgressDefinition[T]) setTick(b bool) {
	d.tick = b || d.tick
}

func (d *ProgressDefinition[T]) SetAutoClose(b ...bool) T {
	d.autoclose = optionutils.BoolOption(b...)
	return d.Self()
}

func (d *ProgressDefinition[T]) IsAutoClose() bool {
	return d.autoclose
}

// SetDecoratorFormat sets the output format for the next decorator.
func (d *ProgressDefinition[T]) SetDecoratorFormat(f ...ttycolors.FormatProvider) T {
	d.nextdecoratorFormat = ttycolors.New(f...)
	return d.Self()
}

func (d *ProgressDefinition[T]) SetMinVisualizationColumn(c int) T {
	d.minColumn = c
	return d.Self()
}

func (d *ProgressDefinition[T]) GetMinVisualizationColumn() int {
	return d.minColumn
}

// SetColor sets the output format for the progress indicator line
func (d *ProgressDefinition[T]) SetColor(f ...ttycolors.FormatProvider) T {
	d.format = ttycolors.New(f...)
	return d.Self()
}

func (d *ProgressDefinition[T]) GetColor() ttycolors.Format {
	return d.format
}

// SetProgressColor sets the output format for the progress indicator.
func (d *ProgressDefinition[T]) SetProgressColor(f ...ttycolors.FormatProvider) T {
	d.progressFormat = ttycolors.New(f...)
	return d.Self()
}

func (d *ProgressDefinition[T]) GetProgressColor() ttycolors.Format {
	return d.progressFormat
}

func format(fmt *ttycolors.Format, def DecoratorDefinition) DecoratorDefinition {
	if *fmt == nil {
		return def
	}
	eff := *fmt
	*fmt = nil
	return Formatted(def, eff)
}

func (d *ProgressDefinition[T]) add(list *[]DecoratorDefinition, def DecoratorDefinition, offset ...int) T {
	if len(offset) == 0 {
		*list = append(*list, format(&d.nextdecoratorFormat, def))
	} else {
		*list = slices.Insert(*list, offset[0], format(&d.nextdecoratorFormat, def))
	}
	return d.Self()
}

// AppendFunc2 runs the decorator function and renders the output on the right of the progress indicator
func (d *ProgressDefinition[T]) AppendFunc2(f DecoratorFunc, offset ...int) {
	d.add(&d.appendDefs, f, offset...)
}

// AppendFunc runs the decorator function and renders the output on the right of the progress indicator
func (d *ProgressDefinition[T]) AppendFunc(f DecoratorFunc, offset ...int) T {
	d.AppendFunc2(f, offset...)
	return d.Self()
}

// AppendDecorator runs the decorator and renders the output on the right of the progress indicator
func (d *ProgressDefinition[T]) AppendDecorator(def DecoratorDefinition, offset ...int) T {
	return d.add(&d.appendDefs, def, offset...)
}

// AppendDecorator2 runs the decorator and renders the output on the right of the progress indicator
func (d *ProgressDefinition[T]) AppendDecorator2(def DecoratorDefinition, offset ...int) {
	d.add(&d.appendDefs, def, offset...)
}

func (d *ProgressDefinition[T]) GetAppendDecorators() []DecoratorDefinition {
	return slices.Clone(d.appendDefs)
}

// PrependFunc2 runs decorator function and render the output left the progress indicator
func (d *ProgressDefinition[T]) PrependFunc2(f DecoratorFunc, offset ...int) {
	d.add(&d.prependDefs, f, offset...)
}

// PrependFunc runs decorator function and render the output left the progress indicator
func (d *ProgressDefinition[T]) PrependFunc(f DecoratorFunc, offset ...int) T {
	d.PrependFunc2(f, offset...)
	return d.Self()
}

// PrependDefinition runs the decorator and renders the output on the right of the progress indicator
func (d *ProgressDefinition[T]) PrependDecorator(def DecoratorDefinition, offset ...int) T {
	return d.add(&d.prependDefs, def, offset...)
}

// PrependDefinition2 runs the decorator and renders the output on the right of the progress indicator
func (d *ProgressDefinition[T]) PrependDecorator2(def DecoratorDefinition, offset ...int) {
	d.add(&d.prependDefs, def, offset...)
}

func (d *ProgressDefinition[T]) GetPrependDecorators() []DecoratorDefinition {
	return slices.Clone(d.prependDefs)
}

// AppendElapsed appends the time elapsed to the progress indicator
func (d *ProgressDefinition[T]) AppendElapsed(offset ...int) T {
	d.tick = true
	return d.AppendFunc(timeElapsed, offset...)
}

// PrependElapsed prepends the time elapsed to the beginning of the indicator
func (d *ProgressDefinition[T]) PrependElapsed(offset ...int) T {
	d.tick = true
	return d.PrependFunc(timeElapsed, offset...)
}

// AppendMessage appends text to the progress indicator
func (d *ProgressDefinition[T]) AppendMessage(m string, offset ...int) T {
	return d.AppendFunc(Message(m), offset...)
}

// PrependMessage prepends text to the beginning of the indicator
func (d *ProgressDefinition[T]) PrependMessage(m string, offset ...int) T {
	return d.PrependFunc(Message(m), offset...)
}

// AppendMVariable appends text to the progress indicator
func (d *ProgressDefinition[T]) AppendVariable(m string, offset ...int) T {
	return d.AppendFunc(Variable(m), offset...)
}

// PrependVariable prepends text to the beginning of the indicator
func (d *ProgressDefinition[T]) PrependVariable(m string, offset ...int) T {
	return d.PrependFunc(Variable(m), offset...)
}

func timeElapsed(e ElementState) any {
	var s string
	if e.IsStarted() {
		s = PrettyTime(e.TimeElapsed())
	}
	return stringutils.PadLeft(s, 5, ' ')
}

////////////////////////////////////////////////////////////////////////////////

// ProgressSpecification is the configuration interface for progress indicators.
type ProgressSpecification[T any] interface {
	ElementSpecification[T]

	// SetAutoClose enables/disabled automatic closeing
	// the element when Update indicates finished.
	SetAutoClose(b ...bool) T

	// SetColor set the color used for the progress line.
	SetColor(col ...ttycolors.FormatProvider) T

	// SetProgressColor set the color used for the progress visualization.
	SetProgressColor(col ...ttycolors.FormatProvider) T

	// SetDecoratorFormat set the output format for the next decorator.
	SetDecoratorFormat(col ...ttycolors.FormatProvider) T

	// SetMinVisualizationColumn sets the minimal column for the visualization,
	SetMinVisualizationColumn(int) T

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

	// AppendDecorator adds a definition providing some text appended
	// to the basic progress indicator.
	// If there are implicit settings, the offset can be used to
	// specify the index in the list of functions.
	AppendDecorator(f DecoratorDefinition, offset ...int) T

	// PrependDecorator adds a definition providing some text prepended
	// to the basic progress indicator.
	// If there are implicit settings, the offset can be used to
	// specify the index in the list of functions.
	PrependDecorator(f DecoratorDefinition, offset ...int) T

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

	setTick(b bool)
}

// ProgressConfiguration provides access to the generic Progress configuration
type ProgressConfiguration interface {
	ElementConfiguration
	GetTick() bool
	IsAutoClose() bool

	GetColor() ttycolors.Format
	GetProgressColor() ttycolors.Format
	GetPrependDecorators() []DecoratorDefinition
	GetAppendDecorators() []DecoratorDefinition
	GetMinVisualizationColumn() int
}

////////////////////////////////////////////////////////////////////////////////

func TransferProgressConfig[D ProgressSpecification[T], T any](d D, c ProgressConfiguration) D {
	for _, e := range c.GetPrependDecorators() {
		d.PrependDecorator(e)
	}
	for _, e := range c.GetAppendDecorators() {
		d.AppendDecorator(e)
	}
	d.SetAutoClose(c.IsAutoClose())
	d.SetColor(c.GetColor())
	d.setTick(c.GetTick())
	return TransferElementConfig(d, c)
}

type Prepender interface {
	PrependFunc2(f DecoratorFunc, offset ...int)
	PrependDecorator2(f DecoratorDefinition, offset ...int)
}

type Appender interface {
	AppendFunc2(f DecoratorFunc, offset ...int)
	AppendDecorator2(f DecoratorDefinition, offset ...int)
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
