package specs

import (
	"github.com/mandelsoft/ttyprogress/types"
)

type GroupInterface interface {
	Container
	ElementInterface
}

// GroupNotifier is used to propagate
// group changes to the main group progress indicator.
type GroupNotifier[E ProgressInterface] interface {
	Add(e E, o any)
	Done(e E, o any)
}

// GroupProgressElementDefinition the interface for a progress indicator
// definition to be usable for a main group progress indicator.
type GroupProgressElementDefinition[E ProgressInterface] interface {
	types.ElementDefinition[E]
	GetGroupNotifier() GroupNotifier[E]
}

type DummyGroupNotifier[E any] struct{}

var _ GroupNotifier[ProgressInterface] = (*DummyGroupNotifier[ProgressInterface])(nil)

func (d *DummyGroupNotifier[E]) Add(e E, o any)  {}
func (d *DummyGroupNotifier[E]) Done(e E, o any) {}

type GroupDefinition[T any, E ElementInterface] struct {
	GroupBaseDefinition[T]
	main GroupProgressElementDefinition[E]
}

func NewGroupDefinition[T any, E ElementInterface](s Self[T], main GroupProgressElementDefinition[E]) *GroupDefinition[T, E] {
	return &GroupDefinition[T, E]{
		main:                main,
		GroupBaseDefinition: NewGroupBaseDefinition(s),
	}
}

func (d *GroupDefinition[T, E]) Dup(s Self[T]) GroupDefinition[T, E] {
	dup := *d
	dup.self = s
	return dup
}

func (d *GroupDefinition[T, E]) GetProgress() GroupProgressElementDefinition[E] {
	return d.main
}

////////////////////////////////////////////////////////////////////////////////

type GroupSpecification[T any] interface {
	GroupBaseSpecification[T]
}

type GroupConfiguration[E ElementInterface] interface {
	GroupBaseConfiguration

	// GetProgress provides the main group progress indicator.
	GetProgress() GroupProgressElementDefinition[E]
}
