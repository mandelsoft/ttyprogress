package ppi

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/ttyprogress/blocks"
	"github.com/mandelsoft/ttyprogress/specs"
)

// ElementInterface in the public interface of elements.
type ElementInterface = specs.ElementInterface

type ElementProtected[I any] interface {
	Protected[I]
	Update() bool
}

func ElementSelf[I ElementInterface, P ElementProtected[I]](p P) Self[I, P] {
	return NewSelf[I, P](p)
}

type (
	TitleFormatProvider = specs.TitleFormatProvider
	ViewFormatProvider  = specs.ViewFormatProvider
	TitleLineProvider   = specs.TitleLineProvider
	GapProvider         = specs.GapProvider
	FollowupGapProvider = specs.FollowupGapProvider
)

type ElemBase[I ElementInterface, P ElementProtected[I]] struct {
	Lock sync.RWMutex

	self Self[I, P]

	block  *blocks.Block
	closer func()

	// timeStarted is time progress began.
	timeStarted time.Time
	timeElapsed time.Duration

	closed bool
}

func NewElemBase[I ElementInterface, P ElementProtected[I]](self Self[I, P], p Container, c specs.ElementConfiguration, view int, closer ...func()) (*ElemBase[I, P], error) {
	if view <= 0 {
		view = 1
	}
	b := blocks.NewBlock(view)
	e := &ElemBase[I, P]{self: self, block: b, closer: general.Optional(closer...)}

	b.SetPayload(self.Self())
	if c.GetFinal() != "" {
		b.SetFinal(c.GetFinal())
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
		return nil, err
	}
	return e, nil
}

func (b *ElemBase[I, P]) SetFinal(m string) I {
	b.block.SetFinal(m)
	return b.self.Self()
}

func (b *ElemBase[I, P]) UIBlock() *blocks.Block {
	return b.block
}

func (b *ElemBase[I, P]) Start() {
	if b.start() {
		b.self.Protected().Update()
		b.block.Flush()
	}
}

func (b *ElemBase[I, P]) start() bool {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	var t time.Time
	if b.closed || b.timeStarted != t {
		return false
	}
	b.timeStarted = time.Now()
	return true
}

func (b *ElemBase[I, P]) IsStarted() bool {
	b.Lock.RLock()
	defer b.Lock.RUnlock()

	var t time.Time
	return b.timeStarted != t
}

func (b *ElemBase[I, P]) Close() error {
	err := b.close()

	if err == nil {
		if b.closer != nil {
			b.closer()
		}
		b.self.Protected().Update()
		b.block.Close()
	}
	return err
}

func (b *ElemBase[I, P]) close() error {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	if b.closed {
		return os.ErrClosed
	}
	b.closed = true
	b.timeElapsed = time.Since(b.timeStarted)
	return nil
}

func (b *ElemBase[I, P]) IsClosed() bool {
	b.Lock.RLock()
	defer b.Lock.RUnlock()

	return b.closed
}

func (b *ElemBase[I, P]) Wait(ctx context.Context) error {
	return b.block.Wait(ctx)
}

// TimeElapsed returns the time elapsed
func (b *ElemBase[I, P]) TimeElapsed() time.Duration {
	if !b.IsStarted() {
		return 0
	}
	b.Lock.RLock()
	defer b.Lock.RUnlock()

	if b.closed {
		return b.timeElapsed
	}
	return time.Since(b.timeStarted)
}

// TimeElapsedString returns the formatted string representation of the time elapsed
func (b *ElemBase[I, P]) TimeElapsedString() string {
	if b.IsStarted() {
		return specs.PrettyTime(b.TimeElapsed())
	}
	return ""
}
