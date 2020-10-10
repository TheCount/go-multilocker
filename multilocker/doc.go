// Package multilocker allows you to atomically lock multiple sync.Lockers at
// once while avoiding deadlocks.
//
// Create a new multilocker with the New function. The returned multilocker
// satisfies the sync.Locker interface.
//
// There are some rules for the safe use of multilockers, most of which are
// similar to the rules concering single lockers:
//
//     * Multilockers are not recursive. Do not try to lock a locked
//       multilocker. Do not try to unlock an unlocked multilocker.
//     * Do not copy a multilocker.
//     * Do not try to create a multilocker from both the write and the read
//       locker of a sync.RWMutex.
//     * Do not make a multilocker part of another multilocker.
//
// Deadlock safety
//
// To use multilockers in a deadlock-safe manner, there is a simple rule: each
// goroutine should not hold more than one lock (on a multilocker or other
// locker) at a time.
//
// Limitations
//
// There is no optimal way to avoid deadlocks, and this package, too, has some
// important limitations:
//
//     * Since you should hold at most one lock per goroutine, you need to know
//       in advance which lockers to group into a multilocker. This is not
//       always easily possible.
//     * Multilockers do not prevent deadlocks which happen in the context of
//       intrinsic locking in go, such as with channel operations or certain
//       library functions.
//
// Finally, remember that use of sync.Lockers (and thus multilockers) should be
// an exception rather than the norm. Share memory by communicating, do not
// communicate by sharing memory. While go channels are not intrinsically
// deadlock-safe, it is generally easier to design a deadlock-free
// communication pattern with channels than with sync.Lockers.
//
// Efficiency
//
// Unlocking a multilocker takes O(n) time, where n is the number of underlying
// locks. The same goes for locking, not counting the time Lock is blocked
// on an underlying lock held by a different goroutine, of course.
//
// Creating a new multilocker from n lockers takes O(n*log(n)) amortised time.
package multilocker
