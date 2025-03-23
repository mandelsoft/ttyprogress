package specs

import (
	"github.com/mandelsoft/ttycolors"
)

type TextSpinnerInterface interface {
	SpinnerInterface
	TextInterface
}

type TextSpinnerDefinition[T any] struct {
	SpinnerDefinition[T]

	view       int
	viewFormat ttycolors.Format
	gap        string
}

var (
	_ TextSpinnerSpecification[any] = (*TextSpinnerDefinition[any])(nil)
	_ TextSpinnerConfiguration      = (*TextSpinnerDefinition[any])(nil)
)

// NewTextSpinnerDefinition can be used to create a nested definition
// for a derived text spinner definition.
func NewTextSpinnerDefinition[T any](self Self[T]) TextSpinnerDefinition[T] {
	d := TextSpinnerDefinition[T]{view: 3}
	d.SpinnerDefinition = NewSpinnerDefinition(self)
	return d
}

func (d *TextSpinnerDefinition[T]) Dup(s Self[T]) TextSpinnerDefinition[T] {
	dup := *d
	dup.SpinnerDefinition = d.SpinnerDefinition.Dup(s)
	return dup
}

func (d *TextSpinnerDefinition[T]) SetView(view int) T {
	d.view = view
	return d.Self()
}

func (d *TextSpinnerDefinition[T]) GetView() int {
	return d.view
}

func (d *TextSpinnerDefinition[T]) SetViewFormat(f ...ttycolors.FormatProvider) T {
	d.viewFormat = ttycolors.New(f...)
	return d.Self()
}

func (d *TextSpinnerDefinition[T]) GetViewFormat() ttycolors.Format {
	return d.viewFormat
}

func (d *TextSpinnerDefinition[T]) SetFollowUpGap(gap string) T {
	d.gap = gap
	return d.Self()
}

func (d *TextSpinnerDefinition[T]) GetFollowUpGap() string {
	return d.gap
}

////////////////////////////////////////////////////////////////////////////////

type TextSpinnerSpecification[T any] interface {
	SpinnerSpecification[T]
	SetView(view int) T
	SetFollowUpGap(gap string) T
	SetViewFormat(f ...ttycolors.FormatProvider) T
}

type TextSpinnerConfiguration interface {
	SpinnerConfiguration
	FollowupGapProvider
	ViewFormatProvider
	GetView() int
}
