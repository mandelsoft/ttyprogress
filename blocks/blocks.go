package blocks

import (
	"context"
	"io"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/mandelsoft/goutils/general"
)

// Blocks is a sequences of Block/s which represent a trailing range of
// lines on a terminal output given by am output steam. The stream is written to
// update the covered terminal lines with the actual context of the included
// Block/s.
// The contents of the Block/s will be flushed on a timed interval or when
// Flush is called.
type Blocks struct {
	lock sync.RWMutex

	// out is the writer to write to
	out       io.Writer
	termWidth int

	overFlowHandled bool

	blocks    []*Block
	lineCount int

	closed bool
	done   chan struct{}

	timer   *time.Timer
	pending bool
}

// New returns a new Blocks with defaults
func New(opt ...io.Writer) *Blocks {
	w := &Blocks{
		out:   general.OptionalDefaulted[io.Writer](os.Stdout, opt...),
		done:  make(chan struct{}),
		timer: time.NewTimer(0),
	}

	termWidth, _ := getTermSize()
	if termWidth != 0 {
		w.termWidth = termWidth
		w.overFlowHandled = true
	}
	go w.listen()
	return w
}

func (w *Blocks) requestFlush() {
	if w.pending {
		return
	}
	w.pending = true
	w.timer.Reset(time.Millisecond * 250)
}

func (w *Blocks) Done() <-chan struct{} {
	return w.done
}

func (w *Blocks) listen() {
	for {
		select {
		case <-w.done:
			return
		case <-w.timer.C:
			w.Flush()
		}
	}
}

// NewBlock returns a new Block assigned to this Blocks.
func (w *Blocks) NewBlock(view ...int) *Block {
	return w.createBlock(nil, 0, view...)
}

// NewAppendedBlock creates a new assigned Block added after the
// given parent block.
func (w *Blocks) NewAppendedBlock(p *Block, view ...int) *Block {
	return w.createBlock(p, 1, view...)
}

// NewInsertedBlock creates a new assigned Block added before the
// given parent block.
func (w *Blocks) NewInsertedBlock(p *Block, view ...int) *Block {
	return w.createBlock(p, 0, view...)
}

func (w *Blocks) createBlock(p *Block, offset int, view ...int) *Block {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.closed {
		return nil
	}
	b := NewBlock(view...)
	w._addBlock(b, p, offset)
	return b
}

// AppendBlock adds an assigned Block after the
// given parent block.
func (w *Blocks) AppendBlock(b *Block, p *Block) error {
	return w.addBlock(b, p, 1)
}

// InsertBlock adds an unassigned Block before the
// given parent block.
func (w *Blocks) InsertBlock(b *Block, p *Block, view ...int) error {
	return w.addBlock(b, p, 0)
}

func (w *Blocks) AddBlock(b *Block) error {
	return w.addBlock(b, nil, 0)
}

func (w *Blocks) addBlock(b *Block, p *Block, offset int) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.closed {
		return nil
	}
	return w._addBlock(b, p, offset)
}

func (w *Blocks) _addBlock(b *Block, p *Block, offset int) error {
	if !b.blocks.CompareAndSwap(nil, w) {
		return ErrAlreadyAssigned
	}

	if p != nil {
		for i := range w.blocks {
			if w.blocks[i] == p {
				w.blocks = append(w.blocks[:i+offset], append([]*Block{b}, w.blocks[i+offset:]...)...)
				return nil
			}
		}
	}
	w.blocks = append(w.blocks, b)
	return nil
}

func (w *Blocks) Blocks() []*Block {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return slices.Clone(w.blocks)
}

func (w *Blocks) TermWidth() int {
	w.lock.RLock()
	defer w.lock.RUnlock()

	return w.termWidth
}

// Close closed the line range.
// No Blocks can be added anymore, and the
// UIBlocls object is done, when all included Block/s
// are closed.
func (w *Blocks) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.closed {
		return os.ErrClosed
	}
	w.closed = true
	return nil
}

// Wait waits until the object and all included
// Block/s are closed.
// If a context.Context is given it returns
// if the context is done, also.
func (w *Blocks) Wait(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.done:
		return nil
	}
}

func (w *Blocks) discardBlock() error {
	discarded := false
	for len(w.blocks) > 0 && w.blocks[0].closed {
		if !discarded {
			w.clearLines()
			discarded = true
		}
		w.blocks[0].emit(true)
		w.blocks = w.blocks[1:]
	}
	if discarded {
		err := w.flush()
		if w.closed && len(w.blocks) == 0 {
			close(w.done)
		}
		return err
	}
	return nil
}

func (w *Blocks) Flush() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.clearLines()
	w.timer.Stop()
	w.pending = false
	return w.flush()
}

func (w *Blocks) flush() error {
	lines := 0
	for _, b := range w.blocks {
		l, err := b.emit(false)
		lines += l
		if err != nil {
			return err
		}
	}
	w.lineCount = lines
	return err
}
