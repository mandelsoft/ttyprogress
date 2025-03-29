package ppi

import (
	"context"
	"os"
	"time"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
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
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.IsStarted()
}

func (b *ElemBase[I]) IsClosed() bool {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.self.IsClosed()
}

func (b *ElemBase[I]) Start() {
	b.elem.Lock()
	defer b.elem.Unlock()

	b.elem.self.Start()
}

func (b *ElemBase[I]) IsFinished() bool {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.self.IsFinished()
}

func (b *ElemBase[I]) Close() error {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.self.Close()
}

func (b *ElemBase[I]) Hide(f ...bool) {
	b.elem.Lock()
	defer b.elem.Unlock()

	b.elem.self.Hide(f...)
}

func (b *ElemBase[I]) SetFinal(m string) {
	b.elem.Lock()
	defer b.elem.Unlock()

	b.elem.self.SetFinal(m)
}

func (b *ElemBase[I]) Flush() error {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.self.Flush()
}

func (b *ElemBase[I]) Wait(ctx context.Context) error {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.self.Wait(ctx)
}

func (b *ElemBase[I]) TimeElapsed() time.Duration {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.self.TimeElapsed()
}

func (b *ElemBase[I]) Update() bool {
	b.elem.Lock()
	defer b.elem.Unlock()

	return b.elem.self.Update()
}

type ElemBaseImpl[I ElementImpl] struct {
	lock synclog.RWMutex
	self I

	block  *blocks.Block
	closer func()

	// timeStarted is time progress began.
	timeStarted time.Time
	timeElapsed time.Duration

	closed bool
}

// var _ ElementImpl = (*ElemBaseImpl[ElementImpl])(nil)

func NewElemBase[I ElementImpl](self I, p Container, c specs.ElementConfiguration, view int, closer ...func()) (*ElemBase[I], *ElemBaseImpl[I], error) {
	if view <= 0 {
		view = 1
	}

	b := blocks.NewBlock(view)

	e := &ElemBaseImpl[I]{lock: synclog.NewRWMutex(generics.TypeOf[I]().String()), self: self, block: b, closer: general.Optional(closer...)}

	b.SetPayload(self)
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

func (b *ElemBaseImpl[I]) Self() I {
	return b.self
}

func (b *ElemBaseImpl[I]) Lock() {
	b.lock.Lock()
}

func (b *ElemBaseImpl[I]) Unlock() {
	b.lock.Unlock()
}

func (b *ElemBaseImpl[I]) RLock() {
	b.lock.RLock()
}

func (b *ElemBaseImpl[I]) RUnlock() {
	b.lock.RUnlock()
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
		b.self.Update()
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
		b.self.Update()
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
	b.self.Update()
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
