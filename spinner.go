package ttyprogress

import (
	"github.com/mandelsoft/ttycolors"
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
	return &specs.DummyGroupNotifier[Spinner]{}
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
	ppi.SpinnerBase[Spinner]
	closed bool
}

type _spinnerProtected struct {
	*_Spinner
}

func (s *_spinnerProtected) Self() Spinner {
	return s._Spinner
}

func (s *_spinnerProtected) Update() bool {
	return s._update()
}

func (s *_spinnerProtected) Visualize() (ttycolors.String, bool) {
	return s._visualize()
}

// newSpinner creates a Spinner with a predefined
// set of spinner phases taken from SpinnerTypes.
func newSpinner(p Container, c specs.SpinnerConfiguration) (Spinner, error) {
	e := &_Spinner{}
	self := ppi.ProgressSelf[Spinner](&_spinnerProtected{e})
	b, err := ppi.NewSpinnerBase[Spinner](self, p, c, 1, nil)
	if err != nil {
		return nil, err
	}
	e.SpinnerBase = *b
	return e, nil
}

func (s *_Spinner) finalize() {
	s._update()
}

func (s *_Spinner) _update() bool {
	return ppi.Update(s.ProgressBase)
}

func (s *_Spinner) _visualize() (ttycolors.String, bool) {
	return ppi.Visualize(&s.SpinnerBase)
}
