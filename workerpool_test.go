package trabago

import (
	"log"
	"sync"
	"testing"
)

func TestDoWork(t *testing.T) {
	poolSize := 10
	wp := New(uint16(poolSize), true, func(v interface{}) interface{} {
		return 1
	})

	var result int
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for r := range wp.Results() {
			r, ok := r.(int)
			if !ok {
				t.Fatalf("expected int result, got %T result", r)
			}
			result += r
			log.Printf("result: %d", result)
		}
		wg.Done()
	}()

	count := poolSize + 1

	for i := 0; i < count; i++ {
		wp.DoWork(i)
	}
	wp.Run()
	wp.Stop()

	wg.Wait()
	if result != count {
		t.Errorf("expected: %d, got: %d", count, result)
	}
}

func workerfunc(v interface{}) interface{} {
	return v
}
