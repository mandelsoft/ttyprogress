package specs

import (
	"io"

	"github.com/mandelsoft/goutils/general"
)

type TextInterface interface {
	ElementInterface
	io.Writer
}

type TextDefinition[T any] struct {
	ElementDefinition[T]
	view int

	auto        bool
	titleline   string
	gap         string
	followupgap string
}

var (
	_ TextSpecification[any] = (*TextDefinition[any])(nil)
	_ TextConfiguration      = (*TextDefinition[any])(nil)
)

// NewTextDefinition can be used to create a nested definition
// for a derived text definition.
func NewTextDefinition[T any](self Self[T]) TextDefinition[T] {
	d := TextDefinition[T]{view: TextView}
	d.ElementDefinition = NewElementDefinition(self)
	return d
}

func (d *TextDefinition[T]) Dup(s Self[T]) TextDefinition[T] {
	dup := *d
	dup.ElementDefinition = d.ElementDefinition.Dup(s)
	return dup
}

func (d *TextDefinition[T]) SetView(v int) T {
	d.view = v
	return d.Self()
}

func (d *TextDefinition[T]) GetView() int {
	return d.view
}

func (d *TextDefinition[T]) GetAuto() bool {
	return d.auto
}

func (d *TextDefinition[T]) SetAuto(b ...bool) T {
	d.auto = general.OptionalDefaultedBool(true, b...)
	return d.Self()
}

func (d *TextDefinition[T]) SetGap(v string) T {
	d.gap = v
	return d.Self()
}

func (d *TextDefinition[T]) GetGap() string {
	return d.gap
}

func (d *TextDefinition[T]) SetFollowUpGap(v string) T {
	d.followupgap = v
	return d.Self()
}

func (d *TextDefinition[T]) GetFollowUpGap() string {
	return d.followupgap
}

func (d *TextDefinition[T]) SetTitleLine(v string) T {
	d.titleline = v
	return d.Self()
}

func (d *TextDefinition[T]) GetTitleLine() string {
	return d.titleline
}

////////////////////////////////////////////////////////////////////////////////

type TextSpecification[T any] interface {
	ElementSpecification[T]
	SetGap(v string) T

	// SetFollowUpGap provides a prefix for the second and following lines.
	// It could be used together with SetTitleLine.
	SetFollowUpGap(v string) T
	SetTitleLine(v string) T
	SetAuto(b ...bool) T
	SetView(int) T
}

type TextConfiguration interface {
	ElementConfiguration
	TitleLineProvider
	GapProvider
	FollowupGapProvider
	GetView() int
	GetAuto() bool
}
