package types

import (
	"github.com/mandelsoft/goutils/generics"
)

type generalize[G, T Element] struct {
	ElementDefinition[T]
}

func (e *generalize[G, T]) Unwrap() any {
	return e.ElementDefinition
}

func (e *generalize[G, T]) Add(p Container) (G, error) {
	var _nil G
	r, err := e.ElementDefinition.Add(p)
	if err != nil {
		return _nil, err
	}
	return generics.Cast[G](r), nil
}

type Unwrapper interface {
	Unwrap() any
}

func Unwrap(e any) any {
	for {
		if u, ok := e.(Unwrapper); ok {
			e = u.Unwrap()
		} else {
			return e
		}
	}
}

func GeneralizeDefinition[G, T Element](d ElementDefinition[T]) ElementDefinition[G] {
	return &generalize[G, T]{d}
}

func GenericDefinition[T Element](d ElementDefinition[T]) ElementDefinition[Element] {
	return &generalize[Element, T]{d}
}
