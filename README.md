# multilocker

[![Documentation](https://godoc.org/github.com/TheCount/go-multilocker/multilocker?status.svg)](https://godoc.org/github.com/TheCount/go-multilocker/multilocker)

## Install

```sh
go get github.com/TheCount/go-multilocker/multilocker
```

## Usage

For the detailed API, see the [Documentation](https://godoc.org/github.com/TheCount/go-multilocker/multilocker).

Example:

```golang
var mtx1, mtx2 sync.Mutex

// goroutine 1
go func() {
  ml := multilocker.New(&mtx1, &mtx2)
  ml.Lock()
  ml.Unlock()
}()

// goroutine 2
go func() {
  ml := multilocker.New(&mtx2, &mtx1)
  ml.Lock()
  ml.Unlock()
}()
```

This example would be prone to deadlocks if goroutine 1 locked first `mtx1`, then `mtx2`, while goroutine 2 attempted to lock `mtx2` first, and then `mtx1`. The use of a multilocker makes locking and unlocking atomic and deadlock-free in this example.
