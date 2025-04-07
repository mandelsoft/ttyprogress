package specs

import (
	"time"
)

const Tick = time.Millisecond * 20

type Speed struct {
	interval time.Duration
	last     time.Time
	passed   time.Duration
}

func (t *Speed) SetSpeed(n int) {
	t.interval = Tick * 5 * time.Duration(n)
}

func (t *Speed) Tick1() bool {
	t.passed += Tick

	if t.interval <= t.passed {
		t.passed -= t.interval
		return true
	}
	return false
}

func (t *Speed) Tick() bool {
	now := time.Now()
	diff := now.Sub(t.last)

	if t.interval <= diff {
		t.last = now
		return true
	}
	return false
}

func NewSpeed(n int) *Speed {
	t := &Speed{}
	t.SetSpeed(n)
	return t
}
