package multilocker

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestNewNoArgs tests calling New without arguments.
func TestNewNoArgs(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic on New() without args")
		}
	}()
	New()
}

// Test a multilocker made of a single locker.
func TestSingle(t *testing.T) {
	var mtx sync.Mutex
	ml := New(&mtx)
	ml.Lock()
	ml.Unlock()
}

// TestTriplet tests multilockers in a triplet configuration.
func TestTriplet(t *testing.T) {
	var mtx1, mtx2, mtx3 sync.Mutex
	var wg sync.WaitGroup
	wg.Add(3)
	testfunc := func(l1, l2 sync.Locker) {
		defer wg.Done()
		ml := New(l1, l2)
		ml.Lock()
		ml.Unlock()
	}
	go testfunc(&mtx1, &mtx2)
	go testfunc(&mtx2, &mtx3)
	go testfunc(&mtx3, &mtx1)
	wg.Wait()
	runtime.GC()
	wg.Add(3)
	go testfunc(&mtx2, &mtx1)
	go testfunc(&mtx3, &mtx2)
	go testfunc(&mtx1, &mtx3)
	wg.Wait()
	runtime.GC()
}

// TestHogwild tests many multilockers at once.
func TestHogwild(t *testing.T) {
	if testing.Short() {
		return
	}
	var mtx [100]sync.Mutex
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i != 100; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			perm := rand.Perm(100)
			allLockers := make([]sync.Locker, 100)
			for i := range allLockers {
				allLockers[i] = &mtx[perm[i]]
			}
			lockers := make([]sync.Locker, 0, 100)
			for i := rand.Intn(2); i < 100; i += rand.Intn(2) + 1 {
				lockers = append(lockers, allLockers[i])
			}
			ml := New(lockers...)
			for i := 0; i != 10; i++ {
				ml.Lock()
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				ml.Unlock()
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			runtime.GC()
		}()
	}
	wg.Wait()
}
