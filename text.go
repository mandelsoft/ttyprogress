package ttyprogress

import (
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/object"
	"github.com/mandelsoft/ttyprogress/ppi"
	"github.com/mandelsoft/ttyprogress/specs"
)

// Text provides a range of lines of output
// until the action described by the element is
// calling Text.Close.
// The element can be used as writer to pass the
// intended output.
// After the writer is closed, the complete output is
// shown after earlier elements and outer containers are
// done.
type Text interface {
	specs.TextInterface
}

type TextDefinition struct {
	specs.TextDefinition[*TextDefinition]
}

func NewText(v ...int) *TextDefinition {
	d := &TextDefinition{}
	d.TextDefinition = specs.NewTextDefinition(specs.NewSelf(d))
	lines := general.Optional(v...)
	if lines > 0 {
		d.SetView(lines)
	}
	return d
}

func (d *TextDefinition) Dup() *TextDefinition {
	dup := &TextDefinition{}
	dup.TextDefinition = d.TextDefinition.Dup(specs.NewSelf(dup))
	return dup
}

func (d *TextDefinition) Add(c Container) (Text, error) {
	return newText(c, d)
}

////////////////////////////////////////////////////////////////////////////////

type _Text struct {
	*ppi.ElemBase[*_TextImpl]
	elem *_TextImpl
}

func (t *_Text) Write(data []byte) (int, error) {
	defer t.elem.Lock()()

	return t.elem.Protected().Write(data)
}

type _TextImpl struct {
	*ppi.ElemBaseImpl[*_TextImpl]
}

// NewText creates a new text stream with the given window size.
// With Text.SetAuto updates are triggered by the Text.Write calls.
// Otherwise, Text.Flush must be called to update the text window.
func newText(p Container, c specs.TextConfiguration) (Text, error) {
	e := &_TextImpl{}
	o := &_Text{elem: e}

	b, s, err := ppi.NewElemBase[*_TextImpl](object.NewSelf[*_TextImpl, any](e, o), p, c, c.GetView())
	if err != nil {
		return nil, err
	}
	e.ElemBaseImpl = s
	o.ElemBase = b
	if c.GetAuto() {
		e.Block().SetAuto()
	}
	return o, nil
}

func (t *_TextImpl) Update() bool {
	return false
}

func (t *_TextImpl) Flush() error {
	return t.Block().Flush()
}

func (t *_TextImpl) Write(data []byte) (int, error) {
	return t.Block().Write(data)
}
