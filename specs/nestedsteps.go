package specs

import (
	"slices"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/ttyprogress/types"
)

type NestedStepsInterface interface {
	ProgressInterface
	Incr() (ElementInterface, error)
	Current() ElementInterface
}

type NestedStep struct {
	name string
	def  types.ElementDefinition[ElementInterface]
}

func NewNestedStep[T ElementInterface](name string, def types.ElementDefinition[T]) NestedStep {
	return NestedStep{
		name: name,
		def:  types.GenericDefinition(def),
	}
}

func (n *NestedStep) Name() string {
	return n.name
}

func (n *NestedStep) Definition() types.ElementDefinition[ElementInterface] {
	return n.def
}

type NestedStepsDefinition[T any] struct {
	BarBaseDefinition[T]
	GroupBaseDefinition[T]

	steps         []NestedStep
	showStepTitle bool
}

// NewNestedStepsDefinition can be used to create a nested definition
// for a derived steps definition.
func NewNestedStepsDefinition[T any](self Self[T], steps []NestedStep) NestedStepsDefinition[T] {
	d := NestedStepsDefinition[T]{steps: slices.Clone(steps)}
	d.BarBaseDefinition = NewBarBaseDefinition(self)
	d.GroupBaseDefinition = NewGroupBaseDefinition(self)
	return d

}

func (d *NestedStepsDefinition[T]) Dup(s Self[T]) NestedStepsDefinition[T] {
	dup := *d
	dup.BarBaseDefinition = d.BarBaseDefinition.Dup(s)
	return dup
}

func (d *NestedStepsDefinition[T]) SetSteps(steps []NestedStep) T {
	d.steps = slices.Clone(steps)
	return d.Self()
}

func (d *NestedStepsDefinition[T]) GetSteps() []NestedStep {
	return slices.Clone(d.steps)
}

func (d *NestedStepsDefinition[T]) ShowStepTitle(b ...bool) T {
	d.showStepTitle = general.OptionalDefaultedBool(true, b...)
	return d.Self()
}

func (d *NestedStepsDefinition[T]) IsShowStepTitle() bool {
	return d.showStepTitle
}

func (d *NestedStepsDefinition[T]) GetTotal() int {
	return len(d.steps)
}

////////////////////////////////////////////////////////////////////////////////

type NestedStepsSpecification[T any] interface {
	BarBaseSpecification[T]
	GroupBaseSpecification[T]
	SetSteps(steps []NestedStep) T
	ShowStepTitle(b ...bool) T
}

type NestedStepsConfiguration interface {
	GroupBaseConfiguration
	BarBaseConfiguration
	GetSteps() []NestedStep
	IsShowStepTitle() bool
}
