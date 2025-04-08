package types

import (
	"context"
	"io"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress/blocks"
)

type Ticker interface {
	Tick() bool
}

type String ttycolors.String

// DecoratorFunc is a function that can be prepended and appended to the progress bar
// The returned type may be a string or a ttycolors.String. All other types
// are converted to a string value by calling the method String() or using the Go native
// conversion (print format %v).
type DecoratorFunc func(b ElementState) any

func (d DecoratorFunc) CreateDecorator(b ElementState) Decorator {
	return &decoratorFuncDecorator{b, d}
}

type decoratorFuncDecorator struct {
	b ElementState
	f DecoratorFunc
}

func (d *decoratorFuncDecorator) Decorate() any {
	return d.f(d.b)
}

// Decorator provides a piece of text appended or
// prepended to the main progress visualization.
type Decorator interface {
	Decorate() any
}

type Container interface {
	AddBlock(b *blocks.Block) error
	Wait(ctx context.Context) error
}

// Element is the common interface of all
// elements provided by the ttyprogress package
type Element interface {
	io.Closer

	// Hide hides the element.
	Hide(...bool)

	// SetFinal overrides the completion message.
	SetFinal(m string)

	// Start records the actual start time and
	// starts the element.
	Start()

	// IsStarted reports whether element has been started.
	IsStarted() bool

	// IsClosed reports whether element has already been closed.
	IsClosed() bool

	// IsFinished returns whether the progress is done.
	IsFinished() bool

	// TimeElapsed reports the duration this element has been
	// active (time since Start or between Start and Close).
	TimeElapsed() time.Duration

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

// ElementState is the interface for
// passed to decorators to retrieve information
// about the actual element state.
type ElementState interface {
	IsStarted() bool
	IsFinished() bool
	TimeElapsed() time.Duration
}

// ProgressElement is the interface for a
// progress element (like a Bar or a Spinner).
type ProgressElement interface {
	Element
	SetVariable(name string, value any)
	GetVariable(name string) any
}
