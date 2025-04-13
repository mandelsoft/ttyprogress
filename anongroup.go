package ttyprogress

import (
	"github.com/mandelsoft/ttyprogress/blocks"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

// AnonymousGroup is a plain indicator group
// and itself no indicator. It can be used
// group indicator with common handling
// for gaps and hiding.
type AnonymousGroup interface {
	specs.GroupBaseInterface
	Gap() string
	Close() error
}

type AnonymousGroupDefinition struct {
	specs.GroupBaseDefinition[*AnonymousGroupDefinition]
}

func NewAnonymousGroup() *AnonymousGroupDefinition {
	d := &AnonymousGroupDefinition{}
	d.GroupBaseDefinition = specs.NewGroupBaseDefinition(specs.NewSelf(d))
	return d
}

func (d *AnonymousGroupDefinition) Dup() *AnonymousGroupDefinition {
	dup := &AnonymousGroupDefinition{}
	dup.GroupBaseDefinition = d.GroupBaseDefinition.Dup(specs.NewSelf(d))
	return dup
}

func (d *AnonymousGroupDefinition) Add(c Container) (AnonymousGroup, error) {
	return newAnonymousGroup(c, d)
}

////////////////////////////////////////////////////////////////////////////////

type _anongroup struct {
	*ppi.GroupState
}

var _ AnonymousGroup = (*_anongroup)(nil)

func newAnonymousGroup(p Container, c *AnonymousGroupDefinition) (AnonymousGroup, error) {
	g := &_anongroup{}
	g.GroupState = ppi.NewGroupState(p, c)
	g.AddBlock(blocks.NewBlock(0))
	return g, nil
}
