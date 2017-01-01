package trabago

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
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

	for c := range wp.Callback() {
		log.Print(c)
	}

	if mw.count != count {
		t.Errorf("expected: %d, got: %d", count, mw.count)
	}
}

func TestIsRateLimited(t *testing.T) {
	n := 10
	duration := time.Second

	rl := NewRateLimiter(n, duration)

	for i := 0; i <= n; i++ {
		if rl.IsRateLimited() {
			if i != n {
				t.Errorf("expected: %d, got: %d", n, i)
			}
			continue
		}

		if i == n {
			t.Errorf("expected: anything except %d, got: %d", n, i)
		}
	}
}
