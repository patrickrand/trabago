package trabago

import (
	"sync"
	"time"
)

// A RateLimiter only allows a code block to be executed N times in a give duration of time.
type RateLimiter struct {
	mu       *sync.Mutex
	n        int           // number of times
	duration time.Duration // window of time to buffer on
	buffer   []time.Time   // linked-list would probably be more efficient...
}

func NewRateLimiter(n int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		mu:       new(sync.Mutex),
		n:        n,
		duration: duration,
		buffer:   make([]time.Time, 0),
	}
}

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
