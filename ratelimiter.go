package trabago

import (
	"sync"
	"time"
)

// TODO:
// - buffer should be linked-list, as opposed to a slice
// - create another function, that accepts a func(...interface{}) interface{}, to be executed if not rate-limited

// A RateLimiter only allows a code block to be executed N number of times, in a give duration of time.
type RateLimiter struct {
	mu       *sync.Mutex
	n        int           // number of times
	duration time.Duration // window of time to buffer on
	buffer   []time.Time   // linked-list would probably be more efficient...
}

// NewRateLimiter returns a new rate limiter.
func NewRateLimiter(n int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		mu:       new(sync.Mutex),
		n:        n,
		duration: duration,
		buffer:   make([]time.Time, 0),
	}
}

// IsRateLimited determines whether a code block should be executed, or ignored.
func (rl *RateLimiter) IsRateLimited() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if len(rl.buffer) < rl.n {
		rl.buffer = append(rl.buffer, time.Now())
		return false
	}

	if time.Since(rl.buffer[0]) > rl.duration {
		rl.buffer = append(rl.buffer[1:], time.Now())
		return false
	}

	return true
}
