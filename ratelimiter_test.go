package trabago

import (
	"testing"
	"time"
)

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
