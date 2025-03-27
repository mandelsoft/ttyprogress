package specs

import (
	"github.com/mandelsoft/ttycolors"
)

type SpinnerInterface interface {
	ProgressInterface
}

type SpinnerDefinition[T any] struct {
	ProgressDefinition[T]

	done    string
	speed   int
	phases  Phases
	pending string
}

var (
	_ SpinnerSpecification[ProgressInterface] = (*SpinnerDefinition[ProgressInterface])(nil)
	_ SpinnerConfiguration                    = (*SpinnerDefinition[ProgressInterface])(nil)
)

// NewSpinnerDef can be used to create a nested definition
// for a derived spinner definition.
func NewSpinnerDefinition[T any](self Self[T]) SpinnerDefinition[T] {
	d := SpinnerDefinition[T]{
		ProgressDefinition: NewProgressDefinition(self),
		speed:              SpinnerSpeed,
		done:               Done,
	}
	d.SetPredefined(SpinnerType)
	return d

}

func (d *SpinnerDefinition[T]) Dup(s Self[T]) SpinnerDefinition[T] {
	dup := *d
	dup.ProgressDefinition = d.ProgressDefinition.Dup(s)
	return dup
}

func (d *SpinnerDefinition[T]) SetPredefined(i int) T {
	if c, ok := SpinnerTypes[i]; ok {
		d.phases = NewStaticPhases(c...)
	}
	return d.Self()
}

func (d *SpinnerDefinition[T]) SetDone(m string) T {
	d.done = m
	return d.Self()
}

func (d *SpinnerDefinition[T]) GetDone() string {
	return d.done
}

func (d *SpinnerDefinition[T]) SetPending(m string) T {
	d.pending = m
	return d.Self()
}

func (d *SpinnerDefinition[T]) GetPending() string {
	return d.pending
}

func (d *SpinnerDefinition[T]) SetSpeed(v int) T {
	d.speed = v
	return d.Self()
}

func (d *SpinnerDefinition[T]) GetSpeed() int {
	return d.speed
}

func (d *SpinnerDefinition[T]) SetSimplePhases(p ...string) T {
	d.phases = NewStaticPhases(p...)
	return d.Self()
}

func (d *SpinnerDefinition[T]) SetPhases(p Phases) T {
	d.phases = p
	return d.Self()
}

func (d *SpinnerDefinition[T]) SetFormattedPhases(p ...ttycolors.String) T {
	d.phases = NewStaticFormattedPhases(p...)
	return d.Self()
}

func (d *SpinnerDefinition[T]) GetPhases() Phases {
	return d.phases
}

////////////////////////////////////////////////////////////////////////////////

type SpinnerSpecification[T any] interface {
	ProgressSpecification[T]
	SetPredefined(i int) T
	SetSpeed(v int) T
	SetSimplePhases(p ...string) T
	SetFormattedPhases(p ...ttycolors.String) T
	SetPhases(p Phases) T
	SetDone(string) T
}

type SpinnerConfiguration interface {
	ProgressConfiguration
	GetPending() string
	GetDone() string
	GetSpeed() int
	GetPhases() Phases
}
