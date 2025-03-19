package ppi

import (
	"sync"

	"github.com/mandelsoft/ttyprogress/specs"
)

var SpinnerTypes = specs.SpinnerTypes

type SpinnerBaseInterface interface {
	ProgressInterface
}

type SpinnerBase[P ProgressInterface] struct {
	ProgressBase[P]

	lock sync.Mutex
	self Self[P, ProgressProtected[P]]

	// pending is the message shown before started
	pending string

	// done is the message shown after closed
	done string

	phases []string
	speed  int

	cnt   int
	phase int
}

var _ SpinnerBaseInterface = (*SpinnerBase[ProgressInterface])(nil)

func NewSpinnerBase[T ProgressInterface](self Self[T, ProgressProtected[T]], p Container, c specs.SpinnerConfiguration, view int, closer func()) (*SpinnerBase[T], error) {
	e := &SpinnerBase[T]{
		self:    self,
		phases:  c.GetPhases(),
		cnt:     c.GetSpeed() - 1,
		speed:   c.GetSpeed(),
		done:    c.GetDone(),
		pending: c.GetPending(),
	}
	b, err := NewProgressBase[T](self, p, c, view, closer, true)
	if err != nil {
		return nil, err
	}
	e.ProgressBase = *b
	return e, nil
}

func (s *SpinnerBase[T]) SetSpeed(v int) T {
	s.speed = v
	s.cnt = v - 1
	return s.self.Self()
}

func Visualize[T ProgressInterface](s *SpinnerBase[T]) (string, bool) {
	if s.self.Self().IsClosed() {
		return s.done, true
	}
	if !s.self.Self().IsStarted() {
		return s.pending, false
	}
	return s.phases[s.phase], false
}

func (s *SpinnerBase[T]) Tick() bool {
	if s.self == nil || s.self.Self().IsClosed() {
		return false
	}
	s.lock.Lock()

	s.cnt++
	if s.cnt < s.speed {
		s.lock.Unlock()
		return false
	}
	s.cnt = 0
	s.phase = (s.phase + 1) % len(s.phases)
	s.lock.Unlock()
	return s.self.Protected().Update()
}
