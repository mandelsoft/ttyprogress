package ttyprogress

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/blocks"
	"github.com/mandelsoft/ttyprogress/specs"
)

// Context is a set of lines on a terminal
// used to display some live progress information.
// It can be used to display an arbitrary number of
// progress elements, which are independently.
// Leading elements will leave the text window
// used once they are finished.
type Context interface {
	Container

	// GetTTYContext returns the underlying ttycolors.TTYContext.
	GetTTYContext() ttycolors.TTYContext

	IsColorsEnabled() bool
	EnableColors(b ...bool) Context

	// Blocks returns the underlying
	// blocks.Blocks object used
	// to display the progress elements.
	// It can directly be used in combination
	// with progress elements.
	// But all active blocks will prohibit the
	// progress object to complete.
	Blocks() *blocks.Blocks

	// Done returns the done channel.
	// A Context is done, if it is closed and
	// all progress elements are finished.
	Done() <-chan struct{}

	// Close closes the Context. No more
	// progress elements can be added anymore.
	Close() error

	// Wait until the Context is Done.
	// If a context.Context is given, Wait
	// also returns if the context is canceled.
	Wait(ctx context.Context) error
}

type _progress struct {
	lock   sync.Mutex
	blocks *blocks.Blocks
	ticker *time.Ticker

	elements []Element
	closed   bool
}

var _ Container = (*_progress)(nil)

// For creates a new Context, which manages a terminal line range
// used to indicate progress of some actions.
// This line range is always at the end of the given
// writer, which must refer to a terminal device.
// Context indicators are added by explicitly calling
// the appropriate constructors. They take the Context
// they should be attached to as first argument.
func For(opt ...io.Writer) Context {
	p := &_progress{
		blocks: blocks.New(opt...),
		ticker: time.NewTicker(specs.Tick),
	}
	go p.listen()
	return p
}

func (p *_progress) Blocks() *blocks.Blocks {
	return p.blocks
}

func (p *_progress) GetTTYContext() ttycolors.TTYContext {
	return p.blocks.GetTTYGontext()
}

func (p *_progress) IsColorsEnabled() bool {
	return p.Blocks().IsColorsEnabled()
}

func (p *_progress) EnableColors(b ...bool) Context {
	p.Blocks().EnableColors(b...)
	return p
}

func (p *_progress) AddBlock(b *blocks.Block) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.closed {
		return os.ErrClosed
	}
	return p.blocks.AddBlock(b)
}

func (p *_progress) Done() <-chan struct{} {
	return p.blocks.Done()
}

func (p *_progress) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.blocks.CloseOnDone()
	p.closed = true
	return nil
}

func (p *_progress) IsClosed() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.closed
}

func (p *_progress) Wait(ctx context.Context) error {
	return p.blocks.Wait(ctx)
}

func (p *_progress) listen() {
	for {
		select {
		case <-p.ticker.C:
			p.tick()
		case <-p.Done():
			p.ticker.Stop()
			return
		}
	}
}

func (p *_progress) tick() {
	flush := false
	for _, b := range p.blocks.Blocks() {
		if e, ok := b.Payload().(Ticker); ok {
			flush = e.Tick() || flush
		}
	}
	if flush {
		p.blocks.Flush()
	}
}
