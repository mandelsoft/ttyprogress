package ppi

import (
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/specs"
)

var SpinnerTypes = specs.SpinnerTypes

type SpinnerInterface interface {
	ProgressInterface
}

type SpinnerImpl interface {
	ProgressImpl
}

type SpinnerBase[P SpinnerImpl] struct {
	*ProgressBase[P]
	elem *SpinnerBaseImpl[P]
}

type SpinnerBaseImpl[P SpinnerImpl] struct {
	*ProgressBaseImpl[P]

	// pending is the message shown before started
	pending string

	// done is the message shown after closed
	done string

	speed *specs.Speed

	phases specs.Phases
}

var _ SpinnerInterface = (*SpinnerBaseImpl[SpinnerImpl])(nil)

func NewSpinnerBase[T SpinnerImpl](self T, p Container, c specs.SpinnerConfiguration, view int, closer func()) (*SpinnerBase[T], *SpinnerBaseImpl[T], error) {
	e := &SpinnerBaseImpl[T]{
		phases:  c.GetPhases(),
		done:    c.GetDone(),
		pending: c.GetPending(),
		speed:   specs.NewSpeed(c.GetSpeed()),
	}
	b, s, err := NewProgressBase[T](self, p, c, view, closer, true)
	if err != nil {
		return nil, nil, err
	}
	e.ProgressBaseImpl = s
	return &SpinnerBase[T]{b, e}, e, nil
}

func (s *SpinnerBaseImpl[T]) Visualize() (ttycolors.String, bool) {
	if s.self.IsClosed() {
		return specs.String(s.done), true
	}
	if !s.self.IsStarted() {
		return specs.String(s.pending), false
	}
	return s.phases.Get(), false
}

func (s *SpinnerBaseImpl[T]) Tick() bool {
	if s.self.IsClosed() {
		return false
	}
	if !s.speed.Tick() {
		return false

	}
	s.phases.Incr()
	return s.self.Update()
}
