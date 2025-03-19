package specs

type CompletedPercent interface {
	CompletedPercent() float64
}

type BarBaseInterface[T any] interface {
	ProgressInterface
	CompletedPercent
	Current() T
}

type BarBaseDefinition[T any] struct {
	ProgressDefinition[T]
	width   uint
	pending string
	config  BarConfig
}

var _ BarBaseSpecification[any] = (*BarBaseDefinition[any])(nil)

// NewBarBaseDefinition can be used to create a nested definition
// for a derived bar definition.
func NewBarBaseDefinition[T any](self Self[T]) BarBaseDefinition[T] {
	return BarBaseDefinition[T]{
		ProgressDefinition: NewProgressDefinition(self),
		config:             BarTypes[BarType],
		width:              BarWidth,
		pending:            Pending,
	}
}

func (d *BarBaseDefinition[T]) Dup(s Self[T]) BarBaseDefinition[T] {
	dup := *d
	dup.ProgressDefinition = d.ProgressDefinition.Dup(s)
	return dup
}

// AppendCompleted appends the completion percent to the progress bar
func (d *BarBaseDefinition[T]) AppendCompleted(offset ...int) T {
	d.AppendFunc(func(b ElementInterface) string {
		return PercentString(b.(CompletedPercent).CompletedPercent())
	}, offset...)
	return d.Self()
}

// PrependCompleted prepends the percent completed to the progress bar
func (d *BarBaseDefinition[T]) PrependCompleted(offset ...int) T {
	d.PrependFunc(func(b ElementInterface) string {
		return PercentString(b.(CompletedPercent).CompletedPercent())
	}, offset...)
	return d.Self()
}

func (d *BarBaseDefinition[T]) SetWidth(w uint) T {
	d.width = w
	return d.Self()
}

func (d *BarBaseDefinition[T]) GetWidth() uint {
	return d.width
}

func (d *BarBaseDefinition[T]) SetPending(m string) T {
	d.pending = m
	return d.Self()
}

func (d *BarBaseDefinition[T]) GetPending() string {
	return d.pending
}

func (d *BarBaseDefinition[T]) SetConfig(c BarConfig) T {
	d.config = c
	return d.Self()
}

func (d *BarBaseDefinition[T]) GetConfig() BarConfig {
	return d.config
}

func (d *BarBaseDefinition[T]) SetPredefined(i int) T {
	if c, ok := BarTypes[i]; ok {
		d.config = c
	}
	return d.Self()
}

func (d *BarBaseDefinition[T]) SetBrackets(c Brackets) T {
	d.config = d.config.SetBrackets(c)
	return d.Self()
}

func (d *BarBaseDefinition[T]) SetBracketType(i int) T {
	d.config = d.config.SetBracketType(i)
	return d.Self()
}

func (d *BarBaseDefinition[T]) SetHead(c rune) T {
	d.config.Head = c
	return d.Self()
}

func (d *BarBaseDefinition[T]) SetEmpty(c rune) T {
	d.config.Empty = c
	return d.Self()
}

func (d *BarBaseDefinition[T]) SetFill(c rune) T {
	d.config.Fill = c
	return d.Self()
}

func (d *BarBaseDefinition[T]) SetLeftEnd(c rune) T {
	d.config.LeftEnd = c
	return d.Self()
}

func (d *BarBaseDefinition[T]) SetRightEnd(c rune) T {
	d.config.RightEnd = c
	return d.Self()
}

////////////////////////////////////////////////////////////////////////////////

type BarBaseSpecification[T any] interface {
	ProgressSpecification[T]
	AppendCompleted(offset ...int) T
	PrependCompleted(offset ...int) T

	SetPending(m string) T
	SetWidth(w uint) T
	SetConfig(c BarConfig) T
	SetPredefined(i int) T
	SetBrackets(c Brackets) T
	SetBracketType(i int) T
	SetHead(c rune) T
	SetEmpty(c rune) T
	SetFill(c rune) T
	SetLeftEnd(c rune) T
	SetRightEnd(c rune) T
}

type BarBaseConfiguration interface {
	ProgressConfiguration
	GetConfig() BarConfig
	GetWidth() uint
	GetPending() string
}

////////////////////////////////////////////////////////////////////////////////

func TransferBarBaseConfig[D BarBaseSpecification[T], T any](d D, c BarBaseConfiguration) D {
	d.SetConfig(c.GetConfig())
	d.SetWidth(c.GetWidth())
	d.SetPending(c.GetPending())
	return TransferProgressConfig(d, c)
}
