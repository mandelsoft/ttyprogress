package ttyprogress

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/object"
	"github.com/mandelsoft/ttyprogress/specs"
)

type BarConfig = specs.BarConfig
type Brackets = specs.Brackets

var (
	BarTypes     = specs.BarTypes
	BracketTypes = specs.BracketTypes

	// BarWidth is the default width of the progress bar
	BarWidth = specs.BarWidth

	// ErrMaxCurrentReached is error when trying to set current value that exceeds the total value
	ErrMaxCurrentReached = errors.New("errors: current value is greater total value")
)

// Bar is a progress bar used to visualize the progress of an action in
// relation to a known maximum of required work.
type Bar interface {
	specs.BarInterface
}

type BarDefinition struct {
	specs.BarDefinition[*BarDefinition]
}

var _ specs.GroupProgressElementDefinition[Bar] = (*BarDefinition)(nil)

func NewBar(set ...int) *BarDefinition {
	d := &BarDefinition{}
	d.BarDefinition = specs.NewBarDefinition(specs.NewSelf(d))
	if len(set) > 0 {
		d.SetPredefined(set[0])
	}
	return d
}

func (d *BarDefinition) Dup() *BarDefinition {
	dup := &BarDefinition{}
	dup.BarDefinition = d.BarDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *BarDefinition) GetGroupNotifier() specs.GroupNotifier {
	return &barGroupNotifier{}
}

func (d *BarDefinition) Add(c Container) (Bar, error) {
	return newBar(c, d)
}

func (d *BarDefinition) AddWithTotal(c Container, total int) (Bar, error) {
	return newBar(c, d, total)
}

////////////////////////////////////////////////////////////////////////////////

// newBar returns a new progress bar
func newBar(p Container, c specs.BarConfiguration[int], total ...int) (Bar, error) {
	b, _, err := newIntBar[IntBarImpl](p, c, general.OptionalDefaulted(c.GetTotal(), total...),
		func(e *IntBarBaseImpl[IntBarImpl], o *IntBarBase[IntBarImpl]) object.Self[IntBarImpl, any] {
			return object.NewSelf[IntBarImpl, any](e, o)
		},
	)
	return b, err
}
