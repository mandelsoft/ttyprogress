package blocks

import (
	"context"
	"io"
	"os"
	"slices"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/ttycolors"
)

const MIN_UPDATE_INTERVAL = 10 * time.Millisecond

// Blocks is a sequences of Block/s which represent a trailing range of
// lines on a terminal output given by am output steam. The stream is written to
// update the covered terminal lines with the actual context of the included
// Block/s.
// The contents of the Block/s will be flushed on a timed interval or when
// Flush is called.
type Blocks struct {
	lock sync.RWMutex

	ttyctx ttycolors.TTYContext
	// out is the writer to write to
	out       io.Writer
	termWidth int

	overFlowHandled bool

	blocks    []*Block
	lineCount int

	closeOnDone bool
	closed      bool
	done        chan struct{}
	cancel      context.CancelFunc
	ctx         context.Context
	request     *request
}

// New returns a new Blocks with defaults
func New(opt ...io.Writer) *Blocks {
	w := &Blocks{
		out:     general.OptionalDefaulted[io.Writer](os.Stdout, opt...),
		done:    make(chan struct{}),
		request: newRequest(),
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())

	if f, ok := w.out.(*os.File); ok {
		w.ttyctx = ttycolors.NewContext(ttycolors.IsTerminal(f))
	}
	termWidth, _ := getTermSize()
	if termWidth != 0 {
		w.termWidth = termWidth
		w.overFlowHandled = true
	}
	go w.listen()
	return w
}

func (w *Blocks) EnableColors(b ...bool) {
	w.ttyctx.Enable(b...)
}

func (w *Blocks) IsColorsEnabled() bool {
	return w.ttyctx.IsEnabled()
}

func (w *Blocks) GetTTYGontext() ttycolors.TTYContext {
	if w == nil {
		return ttycolors.NewContext()
	}
	return w.ttyctx
}

func (w *Blocks) requestFlush() {
	w.request.Request()
}

func (w *Blocks) Done() <-chan struct{} {
	return w.done
}

func (w *Blocks) _flush() {
	w.lock.Lock()
	defer w.lock.Unlock()
	/*
		clearLines(w.out, w.lineCount)
		w.flushAll()
	*/
	w.deltaFlush()
}

func (w *Blocks) listen() {
	for {
		err := w.request.Wait(w.ctx)
		w._flush()
		if err != nil {
			close(w.done)
			return
		}
		time.Sleep(MIN_UPDATE_INTERVAL)
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

func (w *Blocks) NoOfBlocks() int {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return len(w.blocks)
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

func (w *Blocks) CloseOnDone() {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.closeOnDone = true
	w.checkDone()
}

// Close closed the line range.
// No Blocks can be added anymore, and the
// Blocks object is done, when all included Block/s
// are closed.
func (w *Blocks) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.closed {
		return os.ErrClosed
	}
	w.closed = true
	w.checkDone()
	return nil
}

func (w *Blocks) checkDone() {
	if (w.closed || w.closeOnDone) && len(w.blocks) == 0 {
		w.closed = true
		w.cancel()
	}
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
			clearLines(w.out, w.lineCount)
			discarded = true
		}
		w.blocks[0].emit(true)
		w.blocks = w.blocks[1:]
	}
	if discarded {
		err := w.flushAll()
		w.checkDone()
		return err
	}
	return nil
}

func (w *Blocks) Flush() error {
	w.requestFlush()
	return nil
}

func (w *Blocks) flushAll() error {
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

func (w *Blocks) deltaFlush() error {
	// determine first delta block
	start := -1
	complete := 0
	lines := 0
	sum := &complete
	for i, b := range w.blocks {
		if start < 0 {
			if b.updated.Swap(false) {
				start = i
				sum = &lines
			}
		}
		(*sum) += b.lastlines
	}
	if start < 0 {
		return nil
	}

	clearLines(w.out, lines)
	for _, b := range w.blocks[start:] {
		l, err := b.emit(false)
		complete += l
		if err != nil {
			return err
		}
	}
	w.lineCount = complete
	return err
}

////////////////////////////////////////////////////////////////////////////////

type request struct {
	lock sync.Mutex

	sema    *semaphore.Weighted
	pending bool
}

func newRequest() *request {
	r := &request{
		sema: semaphore.NewWeighted(1),
	}
	r.sema.Acquire(context.Background(), 1)
	return r
}

func (r *request) Request() {
	r.lock.Lock()
	defer r.lock.Unlock()

	if !r.pending {
		r.pending = true
		r.sema.Release(1)
	}
}

func (r *request) Wait(ctx context.Context) error {
	if err = r.sema.Acquire(ctx, 1); err != nil {
		return err
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.pending = false
	return nil
}
