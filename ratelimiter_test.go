package trabago

import (
	"testing"
	"time"
)

func TestIsRateLimited(t *testing.T) {
	n := 10
	duration := time.Second

	rl := NewRateLimiter(n, duration)

	for i := 0; i <= 2*n; i++ {
		if rl.IsRateLimited() {
			if i < n {
				t.Errorf("expected a number gte to: %d, got: %d", n, i)
			}
		}
	}
}
