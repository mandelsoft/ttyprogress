package ppi

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/mandelsoft/ttyprogress/blocks"
	"github.com/mandelsoft/ttyprogress/specs"
)

type Gapped interface {
	Gap() string
}

type GroupBase[I any, T ProgressInterface] struct {
	self     I
	lock     sync.RWMutex
	parent   Container
	pgap     string
	gap      string
	followup string

	main     T
	notifier specs.GroupNotifier[T]
	blocks   []*blocks.Block
	closed   bool
}

func NewGroupBase[I any, T ProgressInterface](p Container, self I, c specs.GroupBaseConfiguration, main func(base *GroupBase[I, T]) (T, specs.GroupNotifier[T], error)) (*GroupBase[I, T], T) {
	g := &GroupBase[I, T]{
		self:     self,
		parent:   p,
		gap:      c.GetGap(),
		followup: c.GetFollowUpGap(),
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

func (g *GroupBase[I, T]) AddBlock(b *blocks.Block) error {
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

func (g *GroupBase[I, T]) Flush() error {
	return g.main.Flush()
}

func (g *GroupBase[I, T]) Gap() string {
	return g.pgap + g.followup
}

func (g *GroupBase[I, T]) TimeElapsed() time.Duration {
	return g.main.TimeElapsed()
}

func (g *GroupBase[I, T]) TimeElapsedString() string {
	return g.main.TimeElapsedString()
}

func (g *GroupBase[I, T]) Start() {
	g.main.Start()
}

func (g *GroupBase[I, T]) IsStarted() bool {
	return g.main.IsStarted()
}

func (g *GroupBase[I, T]) Close() error {
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
		g.main.Close()
	}()
	return nil
}

func (g *GroupBase[I, T]) IsClosed() bool {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.closed
}

func (g *GroupBase[I, T]) Wait(ctx context.Context) error {
	return g.blocks[0].Wait(ctx)
}
