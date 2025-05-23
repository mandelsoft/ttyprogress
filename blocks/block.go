package blocks

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"runtime"
	atomic2 "sync/atomic"

	"github.com/mandelsoft/goutils/atomic"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/optionutils"
	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttycolors/ansi"
)

const DefaultView = 10

var ErrNotAssigned = errors.New("uiblock not assigned")
var ErrAlreadyAssigned = errors.New("uiblock already assigned")

// ESC is the ASCII code for escape character
const ESC = 27

// ErrClosedPipe is the error returned when trying to writer is not listening
var ErrClosedPipe = errors.New("uilive: read/write on closed pipe")

// FdWriter is a writer with a file descriptor.
type FdWriter interface {
	io.Writer
	Fd() uintptr
}

type Block struct {
	blocks      atomic.Value[*Blocks]
	titleline   string
	titleFormat ttycolors.Format
	viewFormat  ttycolors.Format
	view        int
	payload     any
	next        *Block // consumer linking
	auto        bool
	gap         string
	followupGap string
	contentGap  string

	startline bool

	buf    bytes.Buffer
	closed bool
	done   chan struct{}

	final       []byte
	hideOnClose bool
	hidden      bool

	updated   atomic2.Bool
	lastlines int

	closer []func()
}

type block = Block

// NewBlock provides a new Block not yet assigned to any Blocks
// object. A Block can only be assigned once.
func NewBlock(view ...int) *Block {
	return &Block{
		startline: true,
		view:      general.OptionalDefaulted(DefaultView, view...),
		done:      make(chan struct{})}
}

func (w *Block) lock() func() {
	b := w.blocks.Load()
	if b == nil {
		return func() {}
	}

	b.lock.Lock()
	return b.lock.Unlock
}

func (w *Block) rlock() func() {
	b := w.blocks.Load()
	if b == nil {
		return func() {}
	}

	b.lock.RLock()
	return b.lock.RUnlock
}

func (w *Block) Blocks() *Blocks {
	return w.blocks.Load()
}

func (w *Block) RegisterCloser(f func()) {
	defer w.lock()()
	w.closer = append(w.closer, f)
}

func (w *Block) HideOnClose(b ...bool) *Block {
	w.hideOnClose = optionutils.BoolOption(b...)
	return w
}

func (w *Block) IsHideOnClose() bool {
	defer w.rlock()()
	return w.hideOnClose
}

func (w *Block) Hide(b ...bool) *Block {
	w.hidden = optionutils.BoolOption(b...)
	return w
}

func (w *Block) IsHidden() bool {
	defer w.rlock()()
	return w.hidden
}

func (w *Block) SetTitleFormat(f ttycolors.Format) *Block {
	defer w.lock()()

	w.titleFormat = f
	w.Flush()
	return w
}

func (w *Block) SetViewFormat(f ttycolors.Format) *Block {
	defer w.lock()()

	w.viewFormat = f
	return w
}

func (w *Block) SetTitleLine(s string) *Block {
	defer w.lock()()

	w.titleline = s
	w.Flush()
	return w
}

func (w *Block) SetFinal(data string) *Block {
	defer w.lock()()

	w.final = []byte(data)
	w.Flush()
	return w
}

func (w *Block) SetAuto(b ...bool) *Block {
	defer w.lock()()

	w.auto = general.OptionalDefaultedBool(true, b...)
	w.Flush()
	return w
}

func (w *Block) GetGap() string {
	defer w.rlock()()
	return w.gap
}

func (w *Block) SetGap(gap string) *Block {
	defer w.lock()()

	if w.followupGap == "" {
		w.followupGap = gap
	}
	w.gap = gap
	w.Flush()
	return w
}

func (w *Block) SetFollowUpGap(gap string) *Block {
	defer w.lock()()

	w.followupGap = gap
	w.Flush()
	return w
}

func (w *Block) SetContentGap(gap string) *Block {
	defer w.lock()()

	w.contentGap = gap
	w.Flush()
	return w
}

func (w *Block) SetPayload(p any) *Block {
	defer w.lock()()

	w.payload = p
	return w
}

func (w *Block) Payload() any {
	defer w.rlock()()

	return w.payload
}

// SetNext set consumer linking.
// The methods SetNext and Next are
// intended for consumers to maintain
// own block sequences.
func (w *Block) SetNext(n *Block) {
	defer w.lock()()

	w.next = n
}

// next get comsumer linking.
// The methods SetNext and Next are
// intended for consumers to maintain
// own block sequences.
func (w *Block) Next() *Block {
	defer w.rlock()()
	return w.next
}

func (w *Block) Reset() {
	defer w.lock()()
	if w.closed {
		return
	}
	w.startline = true
	w.buf.Reset()
}

// Write save the contents of buf to the writer b. The only errors returned are ones encountered while writing to the underlying buffer.
func (w *Block) Write(buf []byte) (n int, err error) {
	defer w.lock()()
	if w.closed {
		return 0, os.ErrClosed
	}

	contentgap := w.followupGap + w.contentGap
	gap := contentgap
	if w.buf.Len() == 0 && w.titleline == "" {
		gap = w.gap + w.contentGap
	}
	if gap != "" {
		for _, b := range buf {
			if b == '\n' {
				w.startline = true
				gap = contentgap
			} else {
				if w.startline {
					w.buf.Write([]byte(gap))
				}
				w.startline = false
			}
			w.buf.WriteByte(b)
		}
	} else {
		n, err = w.buf.Write(buf)
	}
	if w.auto {
		w.Flush()
	}
	return n, err
}

func (w *Block) Flush() error {
	w.updated.Store(true)
	b := w.blocks.Load()
	if b == nil {
		return ErrNotAssigned
	}
	b.requestFlush()
	return nil
}

func (w *Block) Close() error {
	b := w.blocks.Load()
	if b != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}
	if w.closed {
		return os.ErrClosed
	}
	if w.hideOnClose {
		w.Flush()
		w.hidden = true
	}
	w.closed = true
	close(w.done)
	for _, c := range w.closer {
		go c()
		runtime.Gosched()
	}
	if b != nil {
		return b.discardBlock()
	}
	return nil
}

func (w *Block) IsClosed() bool {
	defer w.rlock()()
	return w.closed
}

func (w *Block) Wait(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-w.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type lineinfo struct {
	start    int
	implicit int
}

func (w *Block) _formatTitle(v string) string {
	if w.titleFormat == nil {
		return v
	}
	return w.Blocks().GetTTYGontext().StringWith(w.titleFormat, v).String()
}

func (w *Block) _formatView(v []byte) []byte {
	if w.viewFormat == nil {
		return v
	}
	return []byte(w.Blocks().GetTTYGontext().StringWith(w.viewFormat, v).String())
}

func (w *Block) emit(final bool) (int, error) {
	blocks := w.blocks.Load()

	if w.hidden {
		w.lastlines = 0
		return 0, nil
	}
	lines := 0
	titleline := 0
	newline := false
	data := w.buf.Bytes()
	if w.closed && w.final != nil {
		data = []byte(w._formatTitle(string(w.final)))
	} else {
		if w.titleline != "" {
			blocks.out.Write([]byte(w.gap + w._formatTitle(w.titleline) + "\n"))
			titleline = 1
		}
	}
	if len(data) == 0 {
		w.lastlines = titleline
		return titleline, nil
	}

	implicit := 0
	linestart := make([]lineinfo, w.view)

	escapeSequence := 0

	var col int
	start := 0
	// fmt.Fprintf(os.Stderr, "write [%d] %q\n", len(data), string(data))
	for o, b := range string(data) {
		if escapeSequence == 0 {
			escapeSequence = ansi.EscapeLength(data[o:])
		}
		if escapeSequence == 0 && b == '\n' || (blocks.overFlowHandled && col >= blocks.termWidth) {
			if b != '\n' {
				implicit++
			} else {
				linestart[lines%w.view].start = start
				linestart[lines%w.view].implicit = implicit
				start = o + 1
				lines++
			}
			newline = true
			col = 0
		} else {
			// fmt.Fprintf(os.Stderr, "insert linebreak %d\n", col)
			newline = false
			if escapeSequence > 0 {
				escapeSequence--
			} else {
				col++
			}
		}
	}

	if !newline {
		linestart[lines%w.view].start = start
		linestart[lines%w.view].implicit = implicit
		lines++
		data = append(data, '\n')
	}

	if w.view > 1 {
		newline = false
	}

	var err error
	var eff int

	if final || lines <= w.view {
		_, err = blocks.out.Write(w._formatView(data))
		eff = lines + implicit + titleline
		// fmt.Fprintf(os.Stderr, "data: %s\n", string(data))
		// fmt.Fprintf(os.Stderr, "eff %d, lines %d, implicit %d\n", eff, lines, implicit)

	} else {
		index := (lines) % w.view
		start := linestart[index].start
		view := data[start:]
		_, err = blocks.out.Write(w._formatView(view))
		eff = w.view + implicit - linestart[index].implicit + titleline
		// fmt.Fprintf(os.Stderr, "data: %s\n", string(view))
		// fmt.Fprintf(os.Stderr, "eff %d, lines %d, implicit %d\n", eff, lines, implicit)
	}
	w.lastlines = eff
	return eff, err
}
