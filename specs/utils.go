package specs

import (
	"fmt"
	"time"

	"github.com/mandelsoft/ttyprogress/units"
)

// Message provide a DecoratorFunc for a simple text message.
func Message(m string) DecoratorFunc {
	return func(element ElementInterface) any {
		return m
	}
}

// PercentString returns the formatted string representation of the percent value.
func PercentString(p float64) string {
	return fmt.Sprintf("%3.f%%", p)
}

func PrettyTime(t time.Duration) string {
	if t == 0 {
		return ""
	}
	return units.Seconds(int(t.Truncate(time.Second) / time.Second))
}
