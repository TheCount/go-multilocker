package multilocker

import (
	"runtime"
	"sync"
)

// T encapsulates a multilocker.
// A multilocker must not be copied.
type T struct {
	// list is the list of lockers locked by this multilocker.
	// The list must be in deadlock-safe order.
	list []sync.Locker
}

var _ sync.Locker = &T{}

// setFinalizer sets a finalizer on this multilocker to drop the references to
// its underlying lockers, so the GC can clean those up as well.
func (t *T) setFinalizer() {
	runtime.SetFinalizer(t, func(t *T) {
		for _, l := range t.list {
			put(l)
		}
	})
}

// Lock atomically locks all locks comprising this multilocker.
func (t *T) Lock() {
	for _, l := range t.list {
		l.Lock()
	}
}

// Unlock atomically unlocks all locks comprising this multilocker.
func (t *T) Unlock() {
	// unlock in reverse order
	for i := len(t.list) - 1; i >= 0; i-- {
		t.list[i].Unlock()
	}
}

// New creates a new multilocker. Requires at least one argument.
func New(lockers ...sync.Locker) *T {
	if len(lockers) == 0 {
		panic("need at least one argument")
	}
	result := &T{
		list: optimizeSequence(lockers...),
	}
	result.setFinalizer()
	return result
}
