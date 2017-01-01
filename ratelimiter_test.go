package trabago

import (
	"testing"
	"time"
)

func TestIsRateLimited(t *testing.T) {
	n := 100
	duration := time.Second
	rl := NewRateLimiter(n, duration)

	i := 2 * n
	var counter, total int
	for range time.Tick(time.Millisecond) {
		if !rl.IsRateLimited() {
			counter++
		}
		total++
		if i--; i <= 0 {
			break
		}
	}

	if counter != n {
		t.Errorf("expected: %d, got: %d", n, counter)
	}
}
