package specs

import (
	"github.com/mandelsoft/ttyprogress/types"
)

// GroupNotifier is used to propagate
// group changes to the main group progress indicator.
type GroupNotifier interface {
	Add(e ProgressInterface, o any)
	Done(e ProgressInterface, o any)
}

// GroupProgressElementDefinition the interface for a progress indicator
// definition to be usable for a main group progress indicator.
type GroupProgressElementDefinition[E ProgressInterface] interface {
	types.ElementDefinition[E]
	GetGroupNotifier() GroupNotifier
}

type VoidGroupNotifier struct{}

var _ GroupNotifier = (*VoidGroupNotifier)(nil)

func (d *VoidGroupNotifier) Add(e ProgressInterface, o any)  {}
func (d *VoidGroupNotifier) Done(e ProgressInterface, o any) {}

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
