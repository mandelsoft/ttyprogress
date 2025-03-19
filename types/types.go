package types

import (
	"context"
	"io"
	"time"

	"github.com/mandelsoft/ttyprogress/blocks"
)

// DecoratorFunc is a function that can be prepended and appended to the progress bar
type DecoratorFunc func(b Element) string

type Container interface {
	AddBlock(b *blocks.Block) error
	Wait(ctx context.Context) error
}

// Element is the common interface of all
// elements provided by the ttyprogress package
type Element interface {
	io.Closer

	// Start records the actual start time and
	// starts the element.
	Start()

	// IsStarted reports whether element has been started.
	IsStarted() bool

	// IsClosed reports whether element has already been closed.
	IsClosed() bool

	// TimeElapsed reports the duration this element has been
	// active (time since Start or between Start and Close).
	TimeElapsed() time.Duration

	// TimeElapsedString provides a nice string representation for
	// TimeElapsed.
	TimeElapsedString() string

	// Flush emits the current output.
	Flush() error

	// Wait waits until the element is finished.
	Wait(ctx context.Context) error
}

// ElementDefinition is the common interface for a definition object
// creating an element of type T.
type ElementDefinition[T Element] interface {
	Add(Container) (T, error)
}
