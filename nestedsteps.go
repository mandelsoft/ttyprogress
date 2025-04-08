package ttyprogress

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/mandelsoft/goutils/stringutils"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

// NestedSteps can be used to visualize a sequence of steps
// represented by own progress indicators.
type NestedSteps interface {
	specs.NestedStepsInterface
}

type NestedStep = specs.NestedStep

func NewNestedStep[T Element](name string, definition ElementDefinition[T]) NestedStep {
	return specs.NewNestedStep[T](name, definition)
}

type NestedStepsDefinition struct {
	specs.NestedStepsDefinition[*NestedStepsDefinition]
}

func NewNestedSteps(steps ...specs.NestedStep) *NestedStepsDefinition {
	d := &NestedStepsDefinition{}
	d.NestedStepsDefinition = specs.NewNestedStepsDefinition(specs.NewSelf(d), steps)
	return d
}

func (d *NestedStepsDefinition) Dup() *NestedStepsDefinition {
	dup := &NestedStepsDefinition{}
	dup.NestedStepsDefinition = d.NestedStepsDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *NestedStepsDefinition) Add(c Container) (NestedSteps, error) {
	return newNestedSteps(c, d)
}

////////////////////////////////////////////////////////////////////////////////

type _NestedStepsImpl struct {
	lock sync.Mutex

	steps []NestedStep
	names []string

	main  nestedMain
	group *ppi.GroupBase[nestedMain]
	cur   Element
}

type nestedMain interface {
	ppi.ProgressInterface
	Current() int
	Incr() bool
	IsFinished() bool
}

// newNestedSteps provides a group of step related progress indicators for a
// given set of sequential steps.
// If steptitle is set to true, the step name is reported in the group title.
// If NestedSteps.SetFinal is set to the empty string, only the progress of the
// active step is shown.
func newNestedSteps(p Container, c specs.NestedStepsConfiguration) (NestedSteps, error) {
	var err error

	steps := c.GetSteps()
	names := stringutils.AlignLeft(sliceutils.Transform(steps, func(step NestedStep) string { return step.Name() }), ' ')

	n := &_NestedStepsImpl{steps: steps, names: names}
	n.group, n.main = ppi.NewGroupBase[nestedMain](p, c, func(b *ppi.GroupBase[nestedMain]) (nestedMain, specs.GroupNotifier[nestedMain], error) {
		var d nestedMain
		if c.IsShowStepTitle() {
			d, err = specs.TransferBarBaseConfig(NewSteps(names...), c).Add(b)
		} else {
			d, err = specs.TransferBarBaseConfig(NewBar().SetTotal(len(steps)), c).Add(b)
		}
		return d, &specs.VoidGroupNotifier[nestedMain]{}, nil
	})
	return n, err
}

func (n *_NestedStepsImpl) SetFinal(m string) {
	n.main.SetFinal(m)
}

func (n *_NestedStepsImpl) SetVariable(name string, value any) {
	n.main.SetVariable(name, value)
}

func (n *_NestedStepsImpl) GetVariable(name string) any {
	return n.main.GetVariable(name)
}

func (n *_NestedStepsImpl) IsStarted() bool {
	return n.main.IsStarted()
}

func (n *_NestedStepsImpl) IsClosed() bool {
	return n.group.IsClosed()
}

func (n *_NestedStepsImpl) IsFinished() bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.group.IsFinished()
}

func (n *_NestedStepsImpl) TimeElapsed() time.Duration {
	return n.main.TimeElapsed()
}

func (n *_NestedStepsImpl) HideOnClose(b ...bool) {
	n.group.HideOnClose(b...)
}

func (n *_NestedStepsImpl) Hide(b ...bool) {
	n.group.Hide(b...)
}

func (n *_NestedStepsImpl) Flush() error {
	return n.group.Flush()
}

func (n *_NestedStepsImpl) Start() {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.main.IsStarted() {
		return
	}
	n.main.Start()
	n.add()
}

func (n *_NestedStepsImpl) Current() Element {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.cur
}

func (n *_NestedStepsImpl) add() (Element, error) {
	var err error
	cur := n.main.Current()
	step := n.steps[cur]
	def := step.Definition()
	specs.PrependFunc(def, Message(n.names[cur]), 0)
	n.cur, err = AddElement(n.group, def)
	if err != nil {
		n.cur.Start()
	}
	return n.cur, err
}

func (n *_NestedStepsImpl) Incr() (Element, error) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.cur != nil {
		n.cur.Close()
	}
	n.main.Incr()
	if !n.main.IsFinished() {
		return n.add()
	} else {
		n.group.Close()
	}
	return nil, nil
}

func (n *_NestedStepsImpl) Close() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.main.IsClosed() {
		return os.ErrClosed
	}
	if n.cur != nil {
		n.cur.Close()
	}
	n.cur = nil
	n.group.Close()
	return n.main.Close()
}

////////////////////////////////////////////////////////////////////////////////

func (n *_NestedStepsImpl) Wait(ctx context.Context) error {
	return n.group.Wait(ctx)
}
