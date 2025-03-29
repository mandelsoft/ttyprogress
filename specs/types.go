package specs

import (
	"github.com/mandelsoft/ttyprogress/types"
)

type DecoratorFunc = types.DecoratorFunc

type DecoratorDefinition interface {
	CreateDecorator(e ElementState) types.Decorator
}

type Container = types.Container
