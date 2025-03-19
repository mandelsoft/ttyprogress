package specs

type GroupBaseInterface interface {
	Container
}

type GroupBaseDefinition[T any] struct {
	self     Self[T]
	gap      string
	followup string
}

var _ GroupBaseSpecification[any] = (*GroupBaseDefinition[any])(nil)
var _ GroupBaseConfiguration = (*GroupBaseDefinition[any])(nil)

// NewGroupBaseDefinition can be used to create a nested definition
// for a derived group definition.
func NewGroupBaseDefinition[T any](self Self[T]) GroupBaseDefinition[T] {
	return GroupBaseDefinition[T]{
		self:     self,
		gap:      GroupGap,
		followup: GroupFollowUpGap,
	}
}

func (d *GroupBaseDefinition[T]) Dup(s Self[T]) GroupBaseDefinition[T] {
	dup := *d
	dup.self = s
	return dup
}

func (d *GroupBaseDefinition[T]) SetGap(gap string) T {
	d.gap = gap
	return d.self.Self()
}

func (d *GroupBaseDefinition[T]) GetGap() string {
	return d.gap
}

func (d *GroupBaseDefinition[T]) SetFollowUpGap(gap string) T {
	d.followup = gap
	return d.self.Self()
}

func (d *GroupBaseDefinition[T]) GetFollowUpGap() string {
	return d.followup
}

////////////////////////////////////////////////////////////////////////////////

type GroupBaseSpecification[T any] interface {
	SetGap(string) T
	SetFollowUpGap(string) T
}

type GroupBaseConfiguration interface {
	GetFollowUpGap() string
	GetGap() string
}
