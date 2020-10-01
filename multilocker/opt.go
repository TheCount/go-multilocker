package multilocker

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

// orderItem gives each locker an ID and a reference count so a collection of
// lockers can be brought into an order optmimized for deadlock-freedom.
type orderItem struct {
	// id gives a locker a unique ID. This field is read-only after creation.
	id uint32

	// refcount is the number of multilockers currently using the locker
	// associated with this order item.
	refcount uint32
}

var (
	// nextOrderID is the next ID to be used for an order item.
	nextOrderID uint32

	// order defines an order on registered lockers.
	order = map[sync.Locker]*orderItem{}

	// orderLocker protects access to order.
	orderLocker sync.RWMutex
)

// optimizeSequence returns the specified lockers in an order optimized for
// deadlock-free locking.
//
// The following invariant applies: given an optimized sequence of lockers,
// removing some elements without changing the order yields another
// optimized sequence.
//
// The order items corresponding to the returned lockers will have their
// reference count increased by one, so they should immediately be used in
// a new multilocker.
func optimizeSequence(lockers ...sync.Locker) []sync.Locker {
	// increase refcounts
	for _, l := range lockers {
		take(l)
	}
	// Sort lockers in global order
	sort.Slice(lockers, func(i, j int) bool {
		orderLocker.RLock()
		result := order[lockers[i]].id < order[lockers[j]].id
		orderLocker.RUnlock()
		return result
	})
	return lockers
}

// take increases the reference count of the order item associated with the
// given locker by 1. If there is no such item, it will be created.
func take(l sync.Locker) {
	orderLocker.RLock()
	item := order[l]
	if item != nil {
		atomic.AddUint32(&item.refcount, 1)
		orderLocker.RUnlock()
		return
	}
	orderLocker.RUnlock()
	// Try again
	orderLocker.Lock()
	item = order[l]
	if item == nil {
		item = &orderItem{
			id:       nextOrderID,
			refcount: 1,
		}
		order[l] = item
		nextOrderID++
		orderLocker.Unlock()
		return
	}
	item.refcount++
	orderLocker.Unlock()
}

// put decreases the reference count of the order item associated with the given
// locker by 1. If the reference count reaches zero, the item is removed from
// the order.
func put(l sync.Locker) {
	orderLocker.RLock()
	item := order[l]
	if atomic.AddUint32(&item.refcount, math.MaxUint32) != 0 {
		orderLocker.RUnlock()
		return
	}
	orderLocker.RUnlock()
	// Remove unused entry
	orderLocker.Lock()
	if item.refcount == 0 {
		delete(order, l)
	}
	orderLocker.Unlock()
}
