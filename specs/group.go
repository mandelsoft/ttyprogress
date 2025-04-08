package specs

import (
	"github.com/mandelsoft/ttyprogress/types"
)

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

type VoidGroupNotifier[E any] struct{}

var _ GroupNotifier[ProgressInterface] = (*VoidGroupNotifier[ProgressInterface])(nil)

func (d *VoidGroupNotifier[E]) Add(e E, o any)  {}
func (d *VoidGroupNotifier[E]) Done(e E, o any) {}

////////////////////////////////////////////////////////////////////////////////

type GroupInterface interface {
	Container
	ElementInterface
}

type GroupDefinition[T any, E ProgressInterface] struct {
	GroupBaseDefinition[T]
	main GroupProgressElementDefinition[E]
}

func NewGroupDefinition[T any, E ProgressInterface](s Self[T], main GroupProgressElementDefinition[E]) *GroupDefinition[T, E] {
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

type GroupConfiguration[E ProgressInterface] interface {
	GroupBaseConfiguration

	// GetProgress provides the main group progress indicator.
	GetProgress() GroupProgressElementDefinition[E]
}
