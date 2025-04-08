package ttyprogress

import (
	"github.com/mandelsoft/ttyprogress/specs"
)

// ScrollingScrollingSpinner is a progress indicator without information about
// the concrete progress using a scrolling text to indicate the progress.
type ScrollingSpinner interface {
	specs.ScrollingSpinnerInterface
}

type ScrollingSpinnerDefinition struct {
	specs.ScrollingSpinnerDefinition[*ScrollingSpinnerDefinition]
}

func NewScrollingSpinner(text string, length int) *ScrollingSpinnerDefinition {
	d := &ScrollingSpinnerDefinition{}
	d.ScrollingSpinnerDefinition = specs.NewScrollingSpinnerDefinition(specs.NewSelf(d), text, length)
	return d
}

func (d *ScrollingSpinnerDefinition) GetGroupNotifier() specs.GroupNotifier {
	return &specs.VoidGroupNotifier{}
}

func (d *ScrollingSpinnerDefinition) Dup() *ScrollingSpinnerDefinition {
	dup := &ScrollingSpinnerDefinition{}
	dup.ScrollingSpinnerDefinition = d.ScrollingSpinnerDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *ScrollingSpinnerDefinition) Add(c Container) (Spinner, error) {
	s, err := newSpinner(c, d)
	if s != nil {
		s.Flush()
	}
	return s, err
}
