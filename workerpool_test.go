package trabago

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"testing"
)

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
	mw := &mockWorker{
		Mutex:     new(sync.Mutex),
		WaitGroup: new(sync.WaitGroup),
	}

	count := 21
	callback := make(chan interface{}, count)
	wp := New(10, mw.work, callback)

	for i := 0; i < count; i++ {
		wp.DoWork(i)
	}
	wp.Run()

	wp.Stop()

	for c := range callback {
		log.Print(c)
	}

	if mw.count != count {
		t.Errorf("expected: %d, got: %d", count, mw.count)
	}
}
