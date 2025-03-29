package ppi

// Self represents the effective object,
// the extended self passed to some kind of
// base implementations.
// It contains the effective object public
// object and the protected or implementation
// object.
// To support reentrant synchronization
// the implementation object should not use
// locking, this has to be done by the
// public view. It must implement all
// public methods using locking and
// forwarding the locked call to the
// implementation object.
type Self[P, O any] interface {
	Self() O
	Protected() P
}

type self[P, O any] struct {
	self      O
	protected P
}

func (s *self[P, O]) Self() O {
	return s.self
}

func (s *self[P, O]) Protected() P {
	return s.protected
}

func NewSelf[P, O any](p P, o O) Self[P, O] {
	return &self[P, O]{o, p}
}
