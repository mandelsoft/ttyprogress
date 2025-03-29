package ttyprogress

import (
	"github.com/mandelsoft/goutils/stringutils"
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
	*IntBarBase[*_StepsImpl]
	elem *_StepsImpl
}

func (s *_Steps) GetCurrentStep() string {
	s.elem.Lock()
	defer s.elem.Unlock()

	return s.GetCurrentStep()
}

type _StepsImpl struct {
	*IntBarBaseImpl[*_StepsImpl]
	steps []string
}

// NewSteps create a Steps progress information for a given
// list of sequential steps.
func newSteps(p Container, c specs.StepsConfiguration) (Steps, error) {
	steps := stringutils.AlignLeft(c.GetSteps(), ' ')
	e := &_StepsImpl{steps: steps}

	b, s, err := newIntBar[*_StepsImpl](p, c, len(steps), func(*IntBarBaseImpl[*_StepsImpl]) *_StepsImpl {
		return e
	})
	if err != nil {
		return nil, err
	}
	e.IntBarBaseImpl = s
	return &_Steps{b, e}, nil
}

func (s *_StepsImpl) GetCurrentStep() string {
	c := s.Current()
	if c == 0 && !s.IsStarted() {
		return stringutils.PadRight("", len(s.steps[0]), ' ')
	}
	if c < len(s.steps) {
		return s.steps[c]
	}
	return stringutils.PadRight("", len(s.steps[0]), ' ')
}
