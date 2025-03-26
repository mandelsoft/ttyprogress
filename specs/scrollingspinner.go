package specs

import (
	"slices"
)

type ScrollingSpinnerInterface interface {
	ProgressInterface
}

type ScrollingSpinnerDefinition[T any] struct {
	ProgressDefinition[T]

	done    string
	phases  []string
	pending string
}

var (
	_ ScrollingSpinnerSpecification[ProgressInterface] = (*ScrollingSpinnerDefinition[ProgressInterface])(nil)
	_ ScrollingSpinnerConfiguration                    = (*ScrollingSpinnerDefinition[ProgressInterface])(nil)
)

// NewScrollingSpinnerDef can be used to create a nested definition
// for a derived scrolling spinner definition.
func NewScrollingSpinnerDefinition[T any](self Self[T], text string, length int) ScrollingSpinnerDefinition[T] {
	d := ScrollingSpinnerDefinition[T]{
		ProgressDefinition: NewProgressDefinition(self),
		done:               Done,
	}
	t := text + " "
	if len(t) <= length {
		t = t + " " + text
	}

	s := t + t
	for i := 0; i < len(t); i++ {
		d.phases = append(d.phases, s[i:i+length])
	}
	return d
}

func (d *ScrollingSpinnerDefinition[T]) Dup(s Self[T]) ScrollingSpinnerDefinition[T] {
	dup := *d
	dup.ProgressDefinition = d.ProgressDefinition.Dup(s)
	return dup
}

func (d *ScrollingSpinnerDefinition[T]) SetDone(m string) T {
	d.done = m
	return d.Self()
}

func (d *ScrollingSpinnerDefinition[T]) GetDone() string {
	return d.done
}

func (d *ScrollingSpinnerDefinition[T]) SetPending(m string) T {
	d.pending = m
	return d.Self()
}

func (d *ScrollingSpinnerDefinition[T]) GetPending() string {
	return d.pending
}

func (d *ScrollingSpinnerDefinition[T]) GetSpeed() int {
	return 3
}

func (d *ScrollingSpinnerDefinition[T]) GetPhases() []string {
	return slices.Clone(d.phases)
}

////////////////////////////////////////////////////////////////////////////////

type ScrollingSpinnerSpecification[T any] interface {
	ProgressSpecification[T]
	SetDone(string) T
}

type ScrollingSpinnerConfiguration = SpinnerConfiguration
