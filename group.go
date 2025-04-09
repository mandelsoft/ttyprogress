package ttyprogress

import (
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

type Group interface {
	Container

	Gap() string

	ppi.ProgressInterface
}

type GroupDefinition[E specs.ProgressInterface] struct {
	specs.GroupDefinition[*GroupDefinition[E], E]
}

func NewGroup[E specs.ProgressInterface](p specs.GroupProgressElementDefinition[E]) *GroupDefinition[E] {
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

type _group[E specs.ProgressInterface] struct {
	*ppi.GroupBase[E]
}

var _ Group = (*_group[specs.ProgressInterface])(nil)

func newGroup[E ppi.ProgressInterface](p Container, c specs.GroupConfiguration[E]) (Group, error) {
	g := &_group[E]{}
	g.GroupBase, _ = ppi.NewGroupBase[E](p, c, func(b *ppi.GroupBase[E]) (E, specs.GroupNotifier, error) {
		e, err := c.GetProgress().Add(b)
		if err != nil {
			return e, nil, err
		}
		return e, c.GetProgress().GetGroupNotifier(), err
	})
	return g, nil
}
