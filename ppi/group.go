package ppi

import (
	"context"
	"os"
	"sync"
	atomic2 "sync/atomic"
	"time"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/optionutils"
	"github.com/mandelsoft/ttyprogress/blocks"
	"github.com/mandelsoft/ttyprogress/specs"
)

type Gapped interface {
	Gap() string
}

type GroupState struct {
	lock        sync.RWMutex
	parent      Container
	pgap        string
	gap         string
	followup    string
	hideOnClose bool
	closer      func()

	blocks        []*blocks.Block
	blockinfo     map[*blocks.Block]bool
	notifyCreator func(b *blocks.Block) func()

	closed atomic2.Bool
}

func NewGroupState(p Container, c specs.GroupBaseConfiguration, closer ...func()) *GroupState {
	g := &GroupState{
		parent:      p,
		gap:         c.GetGap(),
		followup:    c.GetFollowUpGap(),
		hideOnClose: c.IsHideOnClose(),
		blocks:      []*blocks.Block{},
		blockinfo:   map[*blocks.Block]bool{},
		closer:      general.Optional(closer...),
	}
	return g
}

// AddBlock adds a group block. There MUST be
// an initial block, typically representing
// the group. It is used as anchor keeping
// the nesting information.
func (g *GroupState) AddBlock(b *blocks.Block) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.closed.Load() {
		return nil
	}

	if len(g.blocks) == 0 {
		err := g.parent.AddBlock(b)
		if err != nil {
			return err
		}
		b.SetGap(g.pgap)
		if g.hideOnClose {
			b.HideOnClose()
		}
	} else {
		n := g.blocks[0]
		for n.Next() != nil && n.Next() != n {
			n = n.Next()
		}
		b.SetGap(g.pgap + g.gap) // .SetFollowUpGap(g.pgap + g.followup)
		g.blocks[0].Blocks().AppendBlock(b, n)
		if g.notifyCreator != nil {
			b.RegisterCloser(g.notifyCreator(b))
		}
		b.RegisterCloser(g.finishBlock)
	}
	if b != nil {
		g.blocks = append(g.blocks, b)
		g.blocks[0].SetNext(b)
	}
	return nil
}

func (g *GroupState) finishBlock() {
	if !g.IsClosed() {
		return
	}
	if g.IsHideOnClose() {
		for _, b := range g.blocks[1:] {
			if !b.IsClosed() {
				return
			}
		}
		for _, b := range g.blocks[1:] {
			b.Hide()
		}
	}
	if g.closer != nil {
		g.closer()
	}
	if !g.blocks[0].IsClosed() {
		g.blocks[0].Close()
	}
}

func (g *GroupState) Gap() string {
	return g.pgap + g.followup
}

func (g *GroupState) HideOnClose(b ...bool) {
	g.blocks[0].HideOnClose(b...)
}

func (g *GroupState) IsHideOnClose() bool {
	if g.hideOnClose {
		return true
	}
	return g.blocks[0].IsHideOnClose()
}

func (g *GroupState) IsHidden() bool {
	return g.blocks[0].IsHidden()
}

func (g *GroupState) Hide(b ...bool) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if optionutils.BoolOption(b...) {
		if g.blocks[0].IsHidden() {
			return
		}
		// save state and hide
		for _, b := range g.blocks[1:] {
			g.blockinfo[b] = b.IsHidden()
			b.Hide()
		}
		g.blocks[0].Hide()
	} else {
		if !g.blocks[0].IsHidden() {
			return
		}
		// restore state
		for _, b := range g.blocks[1:] {
			b.Hide(g.blockinfo[b])
		}
		g.blocks[0].Hide(false)
	}
}

func (g *GroupState) Close() error {
	if g.closed.Swap(true) {
		return os.ErrClosed
	}
	g.finishBlock()
	return nil
}

func (g *GroupState) IsClosed() bool {
	return g.closed.Load()
}

func (g *GroupState) Wait(ctx context.Context) error {
	return g.blocks[0].Wait(ctx)
}

////////////////////////////////////////////////////////////////////////////////

type GroupBase[T ProgressInterface] struct {
	GroupState

	main     T
	notifier specs.GroupNotifier
}

func NewGroupBase[T ProgressInterface](p Container, c specs.GroupBaseConfiguration, main func(base *GroupBase[T]) (T, specs.GroupNotifier, error)) (*GroupBase[T], T) {
	g := &GroupBase[T]{
		GroupState: *NewGroupState(p, c),
	}
	g.notifyCreator = g.notifyCreator
	g.closer = g.closeMain

	if pg, ok := p.(Gapped); ok {
		g.pgap = pg.Gap()
	}

	if m, n, err := main(g); err != nil {
		return nil, g.main
	} else {
		g.main = m
		g.notifier = n
		return g, g.main
	}
}

func (g *GroupBase[T]) AddBlock(b *blocks.Block) error {
	err := g.GroupState.AddBlock(b)
	if err != nil {
		return err
	}
	if len(g.blocks) > 1 {
		g.main.Start()
	}
	return nil
}

func (g *GroupBase[T]) createNotifier(b *blocks.Block) func() {
	g.notifier.Add(g.main, b)
	return func() { g.notifier.Done(g.main, b) }
}

func (g *GroupBase[T]) SetVariable(name string, value any) {
	g.main.SetVariable(name, value)
}

func (g *GroupBase[T]) GetVariable(name string) any {
	return g.main.GetVariable(name)
}

func (g *GroupBase[T]) SetFinal(m string) {
	g.blocks[0].SetFinal(m)
}

func (g *GroupBase[T]) Flush() error {
	return g.main.Flush()
}

func (g *GroupBase[T]) TimeElapsed() time.Duration {
	return g.main.TimeElapsed()
}

func (g *GroupBase[T]) Start() {
	g.main.Start()
}

func (g *GroupBase[T]) IsStarted() bool {
	return g.main.IsStarted()
}

func (g *GroupBase[T]) IsFinished() bool {
	return g.main.IsFinished()
}

func (g *GroupBase[T]) closeMain() {
	g.main.Close()
}
