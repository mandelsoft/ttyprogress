package ttyprogress

// New provides a new independent element definition
// for a given predefined one.
func New[T Dupper[T]](d Dupper[T]) T {
	return d.Dup()
}

// TypeFor provides a new progress indicator
// for a preconfigured archetype,
// A new configuration can be created with
// New. It is preconfigured
// according to the initial configuration and can
// be refined, furthermore.
func TypeFor[T Dupper[T]](d T) Dupper[T] {
	return d.Dup()
}

// GenericDefinition converts a specific progress indicator definition
// into a generic one (providing just an Element instead of the specific type).
func GenericDefinition[T ElementDefinition[E], E Element](d T) ElementDefinition[Element] {
	return &genericDefinition[E]{d}
}

type genericDefinition[E Element] struct {
	def ElementDefinition[E]
}

func (d *genericDefinition[E]) Add(c Container) (Element, error) {
	return d.def.Add(c)
}
