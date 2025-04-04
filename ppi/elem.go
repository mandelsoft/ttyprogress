package ppi

import (
	"context"
	"os"
	"time"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/object"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/blocks"
	"github.com/mandelsoft/ttyprogress/specs"
	"github.com/mandelsoft/ttyprogress/synclog"
	"github.com/mandelsoft/ttyprogress/types"
)

// Element in the public interface of elements.
type Element = types.Element

type ElementImpl interface {
	Element
	/* abstract protected */ Update() bool
}

type (
	TitleFormatProvider = specs.TitleFormatProvider
	ViewFormatProvider  = specs.ViewFormatProvider
	TitleLineProvider   = specs.TitleLineProvider
	GapProvider         = specs.GapProvider
	FollowupGapProvider = specs.FollowupGapProvider
)

type ElemBase[I ElementImpl] struct {
	elem *ElemBaseImpl[I]
}

func (b *ElemBase[I]) IsStarted() bool {
	defer b.elem.Lock()()

	return b.elem.Protected().IsStarted()
}

func (b *ElemBase[I]) IsClosed() bool {
	defer b.elem.Lock()()

	return b.elem.Protected().IsClosed()
}

func (b *ElemBase[I]) Start() {
	defer b.elem.Lock()()

	b.elem.Protected().Start()
}

func (b *ElemBase[I]) IsFinished() bool {
	defer b.elem.Lock()()

	return b.elem.Protected().IsFinished()
}

func (b *ElemBase[I]) Close() error {
	defer b.elem.Lock()()

	return b.elem.Protected().Close()
}

func (b *ElemBase[I]) Hide(f ...bool) {
	defer b.elem.Lock()()

	b.elem.Protected().Hide(f...)
}

func (b *ElemBase[I]) SetFinal(m string) {
	defer b.elem.Lock()()

	b.elem.Protected().SetFinal(m)
}

func (b *ElemBase[I]) Flush() error {
	defer b.elem.Lock()()

	return b.elem.Protected().Flush()
}

func (b *ElemBase[I]) Wait(ctx context.Context) error {
	defer b.elem.Lock()()

	return b.elem.Protected().Wait(ctx)
}

func (b *ElemBase[I]) TimeElapsed() time.Duration {
	defer b.elem.Lock()()

	return b.elem.Protected().TimeElapsed()
}

func (b *ElemBase[I]) Update() bool {
	defer b.elem.Lock()()

	return b.elem.Protected().Update()
}

type ElemBaseImpl[I ElementImpl] struct {
	lock synclog.RWMutex
	self object.Self[I, any]

	block  *blocks.Block
	closer func()

	// timeStarted is time progress began.
	timeStarted time.Time
	timeElapsed time.Duration

	closed bool
}

// var _ ElementImpl = (*ElemBaseImpl[ElementImpl])(nil)

func NewElemBase[I ElementImpl](self object.Self[I, any], p Container, c specs.ElementConfiguration, view int, closer ...func()) (*ElemBase[I], *ElemBaseImpl[I], error) {
	if view <= 0 {
		view = 1
	}

	b := blocks.NewBlock(view)

	e := &ElemBaseImpl[I]{lock: synclog.NewRWMutex(generics.TypeOf[I]().String()), self: self, block: b, closer: general.Optional(closer...)}

	b.SetPayload(self.Self())
	if c.GetFinal() != "" {
		b.SetFinal(c.GetFinal())
	}
	if c.GetHideOnClose() {
		b.HideOnClose(c.GetHideOnClose())
	}
	if c.GetHide() {
		b.Hide(c.GetHide())
	}
	pgap := ""
	if g, ok := p.(Gapped); ok {
		pgap = g.Gap()
	}
	if t, ok := c.(TitleLineProvider); ok && t.GetTitleLine() != "" {
		b.SetTitleLine(t.GetTitleLine())
	}
	if t, ok := c.(GapProvider); ok && t.GetGap() != "" {
		b.SetGap(pgap + t.GetGap())
	}
	if t, ok := c.(FollowupGapProvider); ok && t.GetFollowUpGap() != "" {
		b.SetFollowUpGap(pgap + t.GetFollowUpGap())
	}
	if t, ok := c.(TitleFormatProvider); ok && t.GetTitleFormat() != nil {
		b.SetTitleFormat(t.GetTitleFormat())
	}
	if t, ok := c.(ViewFormatProvider); ok && t.GetViewFormat() != nil {
		b.SetViewFormat(t.GetViewFormat())
	}

	if err := p.AddBlock(b); err != nil {
		return nil, nil, err
	}
	return &ElemBase[I]{e}, e, nil
}

func (b *ElemBaseImpl[I]) Protected() I {
	return b.self.Protected()
}

func (b *ElemBaseImpl[I]) Lock() func() {
	b.lock.Lock()
	return b.lock.Unlock
}

func (b *ElemBaseImpl[I]) RLock() func() {
	b.lock.RLock()
	return b.lock.RUnlock
}

func (b *ElemBaseImpl[I]) HideOnClose(f ...bool) {
	b.block.HideOnClose(f...)
}

func (b *ElemBaseImpl[I]) IsHideOnClose() bool {
	return b.block.IsHideOnClose()
}

func (b *ElemBaseImpl[I]) Hide(f ...bool) {
	b.block.Hide(f...)
}

func (b *ElemBaseImpl[I]) IsHidden() bool {
	return b.block.IsHidden()
}

func (b *ElemBaseImpl[I]) SetFinal(m string) {
	b.block.SetFinal(m)
}

func (b *ElemBaseImpl[I]) Block() *blocks.Block {
	return b.block
}

func (b *ElemBaseImpl[I]) StringWith(f ttycolors.FormatProvider, seq ...any) ttycolors.String {
	return b.block.Blocks().GetTTYGontext().StringWith(f, seq...)
}

func (b *ElemBaseImpl[I]) String(seq ...any) ttycolors.String {
	return b.block.Blocks().GetTTYGontext().String(seq...)
}

func (b *ElemBaseImpl[I]) Start() {
	if b.start() {
		b.Protected().Update()
		b.block.Flush()
	}
}

func (b *ElemBaseImpl[I]) start() bool {
	var t time.Time
	if b.closed || b.timeStarted != t {
		return false
	}
	b.timeStarted = time.Now()
	return true
}

func (b *ElemBaseImpl[I]) IsStarted() bool {
	var t time.Time
	return b.timeStarted != t
}

func (b *ElemBaseImpl[I]) Close() error {
	err := b.close()

	if err == nil {
		if b.closer != nil {
			b.closer()
		}
		b.Protected().Update()
		b.block.Close()
	}
	return err
}

func (b *ElemBaseImpl[I]) close() error {
	if b.closed {
		return os.ErrClosed
	}
	b.closed = true
	b.timeElapsed = time.Since(b.timeStarted)
	return nil
}

func (b *ElemBaseImpl[I]) IsClosed() bool {
	return b.closed
}

func (b *ElemBaseImpl[I]) IsFinished() bool {
	return b.closed
}

func (b *ElemBaseImpl[I]) Flush() error {
	b.Protected().Update()
	return b.block.Flush()
}

func (b *ElemBaseImpl[I]) Wait(ctx context.Context) error {
	return b.block.Wait(ctx)
}

// TimeElapsed returns the time elapsed
func (b *ElemBaseImpl[I]) TimeElapsed() time.Duration {
	if !b.IsStarted() {
		return 0
	}
	if b.closed {
		return b.timeElapsed
	}
	return time.Since(b.timeStarted)
}
