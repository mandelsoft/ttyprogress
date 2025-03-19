package specs

type BarInterface interface {
	BarBaseInterface[int]

	IsFinished() bool
	Set(n int) bool
	Incr() bool
}

type BarDefinition[T any] struct {
	BarBaseDefinition[T]

	total int
}

var (
	_ BarSpecification[any] = (*BarDefinition[any])(nil)
	_ BarConfiguration[int] = (*BarDefinition[any])(nil)
)

// NewBarDefinition can be used to create a nested definition
// for a derived bar definition.
func NewBarDefinition[T any](s Self[T]) BarDefinition[T] {
	return BarDefinition[T]{
		BarBaseDefinition: NewBarBaseDefinition(s),
		total:             100,
	}
}

func (d *BarDefinition[T]) Dup(s Self[T]) BarDefinition[T] {
	dup := *d
	dup.BarBaseDefinition = d.BarBaseDefinition.Dup(s)
	return dup
}

func (d *BarDefinition[T]) SetTotal(v int) T {
	d.total = v
	return d.Self()
}

func (d *BarDefinition[T]) GetTotal() int {
	return d.total
}

////////////////////////////////////////////////////////////////////////////////

type BarSpecification[T any] interface {
	BarBaseSpecification[T]

	SetTotal(v int) T
}

type BarConfiguration[T any] interface {
	BarBaseConfiguration
	GetTotal() int
}
