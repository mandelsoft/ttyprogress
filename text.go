package ttyprogress

import (
	"github.com/mandelsoft/goutils/general"
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
	ppi.ElemBase[Text, *_textProtected]
}

type _textProtected struct {
	*_Text
}

func (t *_textProtected) Self() Text {
	return t._Text
}

func (t *_textProtected) Update() bool {
	return t._update()
}

// NewText creates a new text stream with the given window size.
// With Text.SetAuto updates are triggered by the Text.Write calls.
// Otherwise, Text.Flush must be called to update the text window.
func newText(p Container, c specs.TextConfiguration) (Text, error) {
	e := &_Text{}

	self := ppi.ElementSelf[Text](&_textProtected{e})
	b, err := ppi.NewElemBase[Text](self, p, c, c.GetView())
	if err != nil {
		return nil, err
	}
	e.ElemBase = *b
	if c.GetAuto() {
		b.UIBlock().SetAuto()
	}
	return e, nil
}

func (t *_Text) _update() bool {
	return false
}

func (t *_Text) Flush() error {
	return t.UIBlock().Flush()
}

func (t *_Text) Write(data []byte) (int, error) {
	return t.UIBlock().Write(data)
}
