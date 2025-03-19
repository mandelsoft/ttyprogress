package specs

type Self[T any] interface {
	Self() T
}

type self[T any] struct {
	self T
}

func NewSelf[T any](s T) Self[T] {
	return &self[T]{self: s}
}

func (s *self[T]) Self() T {
	return s.self
}
