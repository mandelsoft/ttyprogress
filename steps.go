package ttyprogress

import (
	"github.com/mandelsoft/goutils/stringutils"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

// Steps can be used to visualize a sequence of steps.
type Steps interface {
	specs.StepsInterface
}

type StepsDefinition struct {
	specs.StepsDefinition[*StepsDefinition]
}

func NewSteps(steps ...string) *StepsDefinition {
	d := &StepsDefinition{}
	d.StepsDefinition = specs.NewStepsDefinition(specs.NewSelf(d), steps)
	return d
}

func (d *StepsDefinition) Dup() *StepsDefinition {
	dup := &StepsDefinition{}
	dup.StepsDefinition = d.StepsDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *StepsDefinition) Add(c Container) (Steps, error) {
	return newSteps(c, d)
}

////////////////////////////////////////////////////////////////////////////////

type _Steps struct {
	*_Bar[Steps]
	steps []string
}

type _stepsProtected struct {
	*_Steps
}

func (b *_stepsProtected) Self() Steps {
	return b._Steps
}

func (b *_stepsProtected) Update() bool {
	return b._update()
}

func (b *_stepsProtected) Visualize() (string, bool) {
	return b._visualize()
}

// NewSteps create a Steps progress information for a given
// list of sequential steps.
func newSteps(p Container, c specs.StepsConfiguration) (Steps, error) {
	steps := stringutils.AlignLeft(c.GetSteps(), ' ')
	e := &_Steps{steps: steps}

	b, err := newBarBase[Steps](p, c, len(steps), func(*_Bar[Steps]) ppi.Self[Steps, ppi.ProgressProtected[Steps]] {
		return ppi.ProgressSelf[Steps](&_stepsProtected{e})
	})

	if err != nil {
		return nil, err
	}
	e._Bar = b
	return e, nil
}

func (s *_Steps) GetCurrentStep() string {
	c := s.Current()
	if c == 0 && !s.IsStarted() {
		return stringutils.PadRight("", len(s.steps[0]), ' ')
	}
	if c < len(s.steps) {
		return s.steps[c]
	}
	return stringutils.PadRight("", len(s.steps[0]), ' ')
}
