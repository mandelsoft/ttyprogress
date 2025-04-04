package ttyprogress

import (
	"github.com/mandelsoft/object"
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
	*ppi.SpinnerBase[*_TextSpinnerImpl]
	elem *_TextSpinnerImpl
}

func (b *_TextSpinner) Write(data []byte) (int, error) {
	defer b.elem.Lock()()

	return b.elem.Protected().Write(data)
}

type _TextSpinnerImpl struct {
	*ppi.SpinnerBaseImpl[*_TextSpinnerImpl]
	closed bool
}

type _textSpinnerProtected struct {
	*_TextSpinnerImpl
}

// newTextSpinner creates  new TextSpinner with the given
// window size. It can be used like a Text element.
func newTextSpinner(p Container, c specs.TextSpinnerConfiguration) (TextSpinner, error) {
	e := &_TextSpinnerImpl{}
	o := &_TextSpinner{elem: e}

	b, s, err := ppi.NewSpinnerBase[*_TextSpinnerImpl](object.NewSelf[*_TextSpinnerImpl, any](e, o), p, c, c.GetView(), nil)
	if err != nil {
		return nil, err
	}
	e.SpinnerBaseImpl = s
	o.SpinnerBase = b
	return o, nil
}

func (b *_TextSpinnerImpl) Write(data []byte) (int, error) {
	b.Start()
	return b.Block().Write(data)
}

func (s *_TextSpinnerImpl) Update() bool {
	line, _ := s.Line()
	s.Block().SetTitleLine(line)
	return true
}
