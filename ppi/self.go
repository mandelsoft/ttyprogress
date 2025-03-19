package ppi

type Protected[I any] interface {
	Self() I
}

// Self represents the effective object,
// the extended self passed to some kind of
// base implementations.
// It contains the effective object
// and a wrapper implementing
// the protected methods required by the
// base implementation but not published
// on the public object interface.
// I is the public effective object interface
// and P the protected implementation wrapper.
type Self[I any, P Protected[I]] interface {
	Self() I
	Protected() P
}

type self[I any, P Protected[I]] struct {
	protected P
}

func (s *self[I, P]) Self() I {
	return s.protected.Self()
}

func (s *self[I, P]) Protected() P {
	return s.protected
}

func NewSelf[I any, P Protected[I]](p P) Self[I, P] {
	return &self[I, P]{p}
}
