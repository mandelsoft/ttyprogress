package ppi

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/mandelsoft/goutils/optionutils"
	"github.com/mandelsoft/ttyprogress/blocks"
	"github.com/mandelsoft/ttyprogress/specs"
)

type Gapped interface {
	Gap() string
}

type GroupBase[T ProgressInterface] struct {
	lock     sync.RWMutex
	parent   Container
	pgap     string
	gap      string
	followup string

	main      T
	notifier  specs.GroupNotifier
	blocks    []*blocks.Block
	blockinfo map[*blocks.Block]bool

	closed bool
}

func NewGroupBase[T ProgressInterface](p Container, c specs.GroupBaseConfiguration, main func(base *GroupBase[T]) (T, specs.GroupNotifier, error)) (*GroupBase[T], T) {
	g := &GroupBase[T]{
		parent:   p,
		gap:      c.GetGap(),
		followup: c.GetFollowUpGap(),
		blocks:   []*blocks.Block{},
	}
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

func (g *GroupBase[T]) SetVariable(name string, value any) {
	g.main.SetVariable(name, value)
}

func (g *GroupBase[T]) GetVariable(name string) any {
	return g.main.GetVariable(name)
}

func (g *GroupBase[T]) SetFinal(m string) {
	g.blocks[0].SetFinal(m)
}

func (g *GroupBase[T]) HideOnClose(b ...bool) {
	g.blocks[0].HideOnClose(b...)
}

func (g *GroupBase[T]) IsHideOnClose() bool {
	return g.blocks[0].IsHideOnClose()
}

func (g *GroupBase[T]) IsHidden() bool {
	return g.blocks[0].IsHidden()
}

func (g *GroupBase[T]) Hide(b ...bool) {
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

func (g *GroupBase[T]) AddBlock(b *blocks.Block) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.closed {
		return nil
	}

	if len(g.blocks) == 0 {
		err := g.parent.AddBlock(b)
		if err != nil {
			return err
		}
		b.SetGap(g.pgap)
	} else {
		g.Start()
		n := g.blocks[0]
		for n.Next() != nil && n.Next() != n {
			n = n.Next()
		}
		b.SetGap(g.pgap + g.gap) // .SetFollowUpGap(g.pgap + g.followup)
		g.blocks[0].Blocks().AppendBlock(b, n)
		b.RegisterCloser(func() { g.notifier.Done(g.main, b) })
		g.notifier.Add(g.main, (b))
	}
	if b != nil {
		g.blocks = append(g.blocks, b)
		g.blocks[0].SetNext(b)
	}
	return nil
}

func (g *GroupBase[T]) Flush() error {
	return g.main.Flush()
}

func (g *GroupBase[T]) Gap() string {
	return g.pgap + g.followup
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

func (g *GroupBase[T]) Close() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.closed {
		return os.ErrClosed
	}
	g.closed = true

	go func() {
		for _, b := range g.blocks[1:] {
			b.Wait(nil)
		}

		if g.blocks[0].IsHideOnClose() {
			for _, b := range g.blocks[1:] {
				b.Hide()
			}
		}
		g.main.Close()
	}()
	return nil
}

func (g *GroupBase[T]) IsClosed() bool {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.closed
}

func (g *GroupBase[T]) IsFinished() bool {
	return g.main.IsFinished()
}

func (g *GroupBase[T]) Wait(ctx context.Context) error {
	return g.blocks[0].Wait(ctx)
}
