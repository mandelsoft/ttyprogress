package synclog

import (
	"sync"
)

type RWMutex interface {
	TryLock() bool
	Lock()
	Unlock()

	TryRLock() bool
	RLock()
	RUnlock()
}

type _RWMutex struct {
	lock sync.RWMutex
	name string
}

var _ RWMutex = (*_RWMutex)(nil)

func NewRWMutex(name string) RWMutex {
	return &_RWMutex{name: name}
}

func (l *_RWMutex) TryLock() bool {
	ok := l.lock.TryLock()
	if ok {
		log("locked rwmutex(%s)[%p](TryLock)\n", l.name, l)
	} else {
		log("not locked rwmutex(%s)[%p](TryLock)\n", l.name, l)
	}
	return ok
}

func (l *_RWMutex) Lock() {
	log("locking rwmutex(%s)[%p]\n", l.name, l)
	l.lock.Lock()
	log("locked rwmutex(%s)[%p]\n", l.name, l)
}

func (l *_RWMutex) Unlock() {
	l.lock.Unlock()
	log("released rwmutex(%s)[%p]\n", l.name, l)
}

func (l *_RWMutex) TryRLock() bool {
	ok := l.lock.TryRLock()
	log("not read locked rwmutex(%s)[%p](TryLock)\n", l.name, l)
	return ok
}

func (l *_RWMutex) RLock() {
	log("read locking rwmutex(%s)[%p]\n", l.name, l)
	l.lock.RLock()
	log("read locked rwmutex(%s)[%p]\n", l.name, l)
}

func (l *_RWMutex) RUnlock() {
	l.lock.RUnlock()
	log("read released rwmutex(%s)[%p]\n", l.name, l)
}
