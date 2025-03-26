package ttyprogress

import (
	"github.com/mandelsoft/ttyprogress/specs"
	"github.com/mandelsoft/ttyprogress/types"
)

// DecoratorFunc is a function that can be prepended and appended to the progress bar
type DecoratorFunc = types.DecoratorFunc

type Element = types.Element

type Container = types.Container

type ElementDefinition[T Element] interface {
	types.ElementDefinition[T]
}

type ElementSpecification[T Element] interface {
	specs.ElementSpecification[T]
}

type Ticker = types.Ticker

type Dupper[T any] interface {
	Dup() T
}
