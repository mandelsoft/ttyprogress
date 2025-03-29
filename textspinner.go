package ttyprogress

import (
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

type TextSpinner interface {
	specs.TextSpinnerInterface
}

type TextSpinnerDefinition struct {
	specs.TextSpinnerDefinition[*TextSpinnerDefinition]
}

func NewTextSpinner(set ...int) *TextSpinnerDefinition {
	d := &TextSpinnerDefinition{}
	d.TextSpinnerDefinition = specs.NewTextSpinnerDefinition(specs.NewSelf(d))
	if len(set) > 0 {
		d.SetPredefined(set[0])
	}
	return d
}

func (d *TextSpinnerDefinition) Dup() *TextSpinnerDefinition {
	dup := &TextSpinnerDefinition{}
	dup.TextSpinnerDefinition = d.TextSpinnerDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *TextSpinnerDefinition) Add(c Container) (TextSpinner, error) {
	return newTextSpinner(c, d)
}

////////////////////////////////////////////////////////////////////////////////

type _TextSpinner struct {
	*ppi.SpinnerBase[TextSpinner]
	closed bool
}

type _textSpinnerProtected struct {
	*_TextSpinner
}

func (t *_textSpinnerProtected) Self() TextSpinner {
	return t._TextSpinner
}

func (t *_textSpinnerProtected) Update() bool {
	return t._update()
}

func (t *_textSpinnerProtected) Visualize() (ttycolors.String, bool) {
	return t._visualize()
}

// newTextSpinner creates  new TextSpinner with the given
// window size. It can be used like a Text element.
func newTextSpinner(p Container, c specs.TextSpinnerConfiguration) (TextSpinner, error) {
	e := &_TextSpinner{}

	self := ppi.ProgressSelf[TextSpinner](&_textSpinnerProtected{e})
	b, err := ppi.NewSpinnerBase[TextSpinner](self, p, c, c.GetView(), nil)
	if err != nil {
		return nil, err
	}
	e.SpinnerBase = b
	return e, nil
}

func (b *_TextSpinner) Write(data []byte) (int, error) {
	b.Start()
	return b.Block().Write(data)
}

func (s *_TextSpinner) _update() bool {
	line, _ := s.Line()
	s.Block().SetTitleLine(line)
	return true
}

func (s *_TextSpinner) _visualize() (ttycolors.String, bool) {
	return ppi.Visualize(s.SpinnerBase)
}
