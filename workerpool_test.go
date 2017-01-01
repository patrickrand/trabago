package trabago

import (
	"errors"
	"strconv"
	"sync"
	"testing"
)

func TestWorkerPool(t *testing.T) {

}

type mockWorker struct {
	*sync.Mutex
	*sync.WaitGroup
	count int
}

func (mw *mockWorker) work(v interface{}) interface{} {
	mw.Lock()
	defer mw.Unlock()
	mw.count++
	if mw.count%2 == 0 {
		return errors.New(strconv.Itoa(mw.count))
	}
	return nil
}

func TestDoWork(t *testing.T) {
	count := 21
	mw := &mockWorker{Mutex: new(sync.Mutex), WaitGroup: new(sync.WaitGroup), count: 0}

	callback := make(chan interface{}, count)

	wp := New(10, mw.work, callback)

	wp.Run()

	for i := 0; i < count; i++ {
		wp.DoWork(i)
	}
	wp.Stop()

	for _ = range wp.Callback() {
	}

	if mw.count != count {
		t.Errorf("expected: %d, got: %d", count, mw.count)
	}
}
