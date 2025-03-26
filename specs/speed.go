package specs

import (
	"time"
)

const Tick = time.Millisecond * 25

type Speed struct {
	interval int64
	passed   int64
}

func (t *Speed) SetSpeed(n int) {
	t.interval = int64(Tick) * 4 * int64(n)
}

func (t *Speed) Tick() bool {
	t.passed += int64(Tick)
	if t.interval <= t.passed {
		t.passed -= t.interval
		return true
	}
	return false
}

func NewSpeed(n int) *Speed {
	t := &Speed{}
	t.SetSpeed(n)
	return t
}
