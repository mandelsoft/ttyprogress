package specs

import (
	"github.com/mandelsoft/goutils/optionutils"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/types"
)

// Element is the common interface of all
// elements provided by the ttyprogress package
type Element = types.Element

type ElementState = types.ElementState

type ElementDefinition[T any] struct {
	self        Self[T]
	final       string
	hideOnClose bool
	hide        bool
}

var (
	_ ElementSpecification[any] = (*ElementDefinition[any])(nil)
	_ ElementConfiguration      = (*ElementDefinition[any])(nil)
)

func NewElementDefinition[T any](self Self[T]) ElementDefinition[T] {
	return ElementDefinition[T]{self: self}
}

func (d *ElementDefinition[T]) Self() T {
	return d.self.Self()
}

func (d *ElementDefinition[T]) Dup(s Self[T]) ElementDefinition[T] {
	dup := *d
	dup.self = s
	return dup
}

func (e *ElementDefinition[T]) HideOnClose(b ...bool) T {
	e.hideOnClose = optionutils.BoolOption(b...)
	return e.self.Self()
}

func (e *ElementDefinition[T]) Hide(b ...bool) T {
	e.hide = optionutils.BoolOption(b...)
	return e.self.Self()
}

func (e *ElementDefinition[T]) GetHide() bool {
	return e.hide
}

func (e *ElementDefinition[T]) GetHideOnClose() bool {
	return e.hideOnClose
}

func (e *ElementDefinition[T]) SetFinal(m string) T {
	e.final = m
	return e.self.Self()
}

func (e *ElementDefinition[T]) GetFinal() string {
	return e.final
}

////////////////////////////////////////////////////////////////////////////////

// TitleLineProvider is the optional interface to provide a title line configuration
// enriching the element configuration.
type TitleLineProvider interface {
	GetTitleLine() string
}

// GapProvider is the optional interface to provide an indent
// enriching the element configuration.
type GapProvider interface {
	GetGap() string
}

// FollowupGapProvider is the optional interface to provide an indent
// for additional lines
// enriching the element configuration.
type FollowupGapProvider interface {
	GetFollowUpGap() string
}

// TitleFormatProvider is the optional interface to provide a
// title format.
type TitleFormatProvider interface {
	GetTitleFormat() ttycolors.Format
}

// ViewFormatProvider is the optional interface to provide a
// view format.
type ViewFormatProvider interface {
	GetViewFormat() ttycolors.Format
}

type ElementSpecification[T any] interface {
	// SetFinal sets a text message shown instead of the
	// text window after the action has been finished.
	SetFinal(string) T
	// HideOnClose will request to hide the element
	// when it is closed.
	HideOnClose(...bool) T

	// Hide will request to initially hide the element.
	Hide(...bool) T
}

type ElementConfiguration interface {
	// GetFinal gets a text message shown instead of the
	// text window after the action has been finished.
	GetFinal() string
	GetHideOnClose() bool
	GetHide() bool
}

////////////////////////////////////////////////////////////////////////////////

func TransferElementConfig[D ElementSpecification[T], T any](d D, c ElementConfiguration) D {
	d.HideOnClose(c.GetHideOnClose())
	d.Hide(c.GetHide())
	d.SetFinal(c.GetFinal())
	return d
}
