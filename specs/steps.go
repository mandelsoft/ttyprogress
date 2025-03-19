package specs

import (
	"slices"
	"time"
)

type StepsInterface interface {
	BarInterface
	GetCurrentStep() string
}

type StepsDefinition[T any] struct {
	BarBaseDefinition[T]

	steps []string
}

// NewStepsDefinition can be used to create a nested definition
// for a derived steps definition.
func NewStepsDefinition[T any](self Self[T], steps []string) StepsDefinition[T] {
	d := StepsDefinition[T]{steps: slices.Clone(steps)}
	d.BarBaseDefinition = NewBarBaseDefinition(self)
	return d

}

func (d *StepsDefinition[T]) Dup(s Self[T]) StepsDefinition[T] {
	dup := *d
	dup.BarBaseDefinition = d.BarBaseDefinition.Dup(s)
	return dup
}

func (d *StepsDefinition[T]) AppendStep() T {
	d.AppendFunc(func(e ElementInterface) string {
		return e.(StepsInterface).GetCurrentStep()
	})
	return d.Self()
}

func (d *StepsDefinition[T]) PrependStep() T {
	d.PrependFunc(func(e ElementInterface) string {
		return e.(StepsInterface).GetCurrentStep()
	})
	return d.Self()
}

func (d *StepsDefinition[T]) SetSteps(steps []string) T {
	d.steps = slices.Clone(steps)
	return d.Self()
}

func (d *StepsDefinition[T]) GetSteps() []string {
	return slices.Clone(d.steps)
}

func (d *StepsDefinition[T]) GetTotal() int {
	return len(d.steps)
}

////////////////////////////////////////////////////////////////////////////////

type StepsSpecification[T any] interface {
	BarBaseSpecification[T]
	SetSteps(steps []string) T
	AppendStep() T
	PrependStep() T
}

type StepsConfiguration interface {
	BarConfiguration[time.Duration]
	GetSteps() []string
}
