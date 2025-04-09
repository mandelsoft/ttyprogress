package ttyprogress

import (
	"github.com/mandelsoft/goutils/stringutils"
	"github.com/mandelsoft/object"
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
	defer s.elem.Lock()()

	return s.elem.Protected().GetCurrentStep()
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
	o := &_Steps{elem: e}

	b, s, err := newIntBar[*_StepsImpl](p, c, len(steps), object.NewSelf[*_StepsImpl, any](e, o))
	if err != nil {
		return nil, err
	}
	e.IntBarBaseImpl = s
	o.IntBarBase = b
	return o, nil
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
