package trabago

import (
	"testing"
	"time"
)

func TestIsRateLimited(t *testing.T) {
	n := 100
	duration := time.Second
	rl := NewRateLimiter(uint(n), duration)

	i := n + (n / 2)
	var counter int
	for range time.Tick(time.Millisecond) {
		if !rl.IsRateLimited() {
			counter++
		}

		if i--; i <= 0 {
			break
		}
	}

	if counter != n {
		t.Errorf("expected: %d, got: %d", n, counter)
	}
}

func TestConditionalExecute(t *testing.T) {
	n := 100
	duration := time.Second
	rl := NewRateLimiter(uint(n), duration)

	i := n + (n / 2)
	var counter int
	for range time.Tick(time.Millisecond) {
		if _, ok := rl.ConditionalExecute(func() interface{} {
			return nil
		}); ok {
			counter++
		}

		if i--; i <= 0 {
			break
		}
	}

	if counter != n {
		t.Errorf("expected: %d, got: %d", n, counter)
	}

}
