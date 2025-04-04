package ttyprogress

import (
	"github.com/mandelsoft/object"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

// Spinner is a progress indicator without information about
// the concrete progress.
type Spinner interface {
	specs.SpinnerInterface
}

type SpinnerDefinition struct {
	specs.SpinnerDefinition[*SpinnerDefinition]
}

func NewSpinner(set ...int) *SpinnerDefinition {
	d := &SpinnerDefinition{}
	d.SpinnerDefinition = specs.NewSpinnerDefinition(specs.NewSelf(d))
	if len(set) > 0 {
		d.SetPredefined(set[0])
	}
	return d
}

func (d *SpinnerDefinition) GetGroupNotifier() specs.GroupNotifier[Spinner] {
	return &specs.VoidGroupNotifier[Spinner]{}
}

func (d *SpinnerDefinition) Dup() *SpinnerDefinition {
	dup := &SpinnerDefinition{}
	dup.SpinnerDefinition = d.SpinnerDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *SpinnerDefinition) Add(c Container) (Spinner, error) {
	return newSpinner(c, d)
}

////////////////////////////////////////////////////////////////////////////////

type _Spinner struct {
	*ppi.SpinnerBase[*_SpinnerImpl]
	elem *_SpinnerImpl
}

type _SpinnerImpl struct {
	*ppi.SpinnerBaseImpl[*_SpinnerImpl]
	closed bool
}

// newSpinner creates a Spinner with a predefined
// set of spinner phases taken from SpinnerTypes.
func newSpinner(p Container, c specs.SpinnerConfiguration) (Spinner, error) {
	e := &_SpinnerImpl{}
	o := &_Spinner{elem: e}
	b, s, err := ppi.NewSpinnerBase[*_SpinnerImpl](object.NewSelf[*_SpinnerImpl, any](e, o), p, c, 1, nil)
	if err != nil {
		return nil, err
	}
	e.SpinnerBaseImpl = s
	o.SpinnerBase = b
	return o, nil
}
