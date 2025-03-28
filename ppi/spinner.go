package ppi

import (
	"sync"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/specs"
)

var SpinnerTypes = specs.SpinnerTypes

type SpinnerBaseInterface interface {
	ProgressInterface
}

type SpinnerBase[P ProgressInterface] struct {
	*ProgressBase[P]

	lock sync.Mutex
	self Self[P, ProgressProtected[P]]

	// pending is the message shown before started
	pending string

	// done is the message shown after closed
	done string

	spped *specs.Speed

	phases specs.Phases
}

var _ SpinnerBaseInterface = (*SpinnerBase[ProgressInterface])(nil)

func NewSpinnerBase[T ProgressInterface](self Self[T, ProgressProtected[T]], p Container, c specs.SpinnerConfiguration, view int, closer func()) (*SpinnerBase[T], error) {
	e := &SpinnerBase[T]{
		self:    self,
		phases:  c.GetPhases(),
		done:    c.GetDone(),
		pending: c.GetPending(),
		spped:   specs.NewSpeed(c.GetSpeed()),
	}
	b, err := NewProgressBase[T](self, p, c, view, closer, true)
	if err != nil {
		return nil, err
	}
	e.ProgressBase = b
	return e, nil
}

func Visualize[T ProgressInterface](s *SpinnerBase[T]) (ttycolors.String, bool) {
	if s.self.Self().IsClosed() {
		return specs.String(s.done), true
	}
	if !s.self.Self().IsStarted() {
		return specs.String(s.pending), false
	}
	return s.phases.Get(), false
}

func (s *SpinnerBase[T]) Tick() bool {
	if s.self == nil || s.self.Self().IsClosed() {
		return false
	}
	s.lock.Lock()
	if !s.spped.Tick() {
		s.lock.Unlock()
		return false

	}
	s.phases.Incr()
	s.lock.Unlock()
	return s.self.Protected().Update()
}
