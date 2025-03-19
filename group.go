package ttyprogress

import (
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
	"github.com/mandelsoft/ttyprogress/types"
)

type Group interface {
	Container

	Gap() string

	ppi.ProgressInterface
}

type GroupDefinition[E types.Element] struct {
	specs.GroupDefinition[*GroupDefinition[E], E]
}

func NewGroup[E types.Element](p specs.GroupProgressElementDefinition[E]) *GroupDefinition[E] {
	d := &GroupDefinition[E]{}
	d.GroupDefinition = *specs.NewGroupDefinition[*GroupDefinition[E], E](specs.NewSelf(d), p)
	return d
}

func (d *GroupDefinition[E]) Dup() *GroupDefinition[E] {
	dup := &GroupDefinition[E]{}
	dup.GroupDefinition = d.GroupDefinition.Dup(specs.NewSelf(d))
	return dup
}

func (d *GroupDefinition[E]) Add(c Container) (Group, error) {
	return newGroup[E](c, d)
}

////////////////////////////////////////////////////////////////////////////////

type _group[E Element] struct {
	*ppi.GroupBase[Group, E]
	main E
}

var _ Group = (*_group[Element])(nil)

func newGroup[E ppi.ProgressInterface](p Container, c specs.GroupConfiguration[E]) (Group, error) {
	g := &_group[E]{}
	g.GroupBase, g.main = ppi.NewGroupBase[Group, E](p, g, c, func(b *ppi.GroupBase[Group, E]) (E, specs.GroupNotifier[E], error) {
		e, err := c.GetProgress().Add(b)
		if err != nil {
			return e, nil, err
		}
		return e, c.GetProgress().GetGroupNotifier(), err
	})
	return g, nil
}
