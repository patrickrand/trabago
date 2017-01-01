package trabago

import (
	"sync"
	"time"
)

type timeNode struct {
	t    time.Time
	next *timeNode
}

// TODO:
// - create another function, that accepts a func(...interface{}) interface{}, to be executed if not rate-limited

// A RateLimiter only allows a code block to be executed N number of times, in a give duration of time.
type RateLimiter struct {
	mu       *sync.Mutex
	n        int           // number of times
	duration time.Duration // window of time to buffer
	bufHead  *timeNode
	bufTail  *timeNode
	bufSize  int
}

// NewRateLimiter returns a new rate limiter.
func NewRateLimiter(n int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		mu:       new(sync.Mutex),
		n:        n,
		duration: duration,
		bufHead:  new(timeNode),
		bufTail:  new(timeNode),
		bufSize:  0,
	}
}

// IsRateLimited determines whether a code block should be executed, or ignored.
func (rl *RateLimiter) IsRateLimited() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// allow execution if buffer is below capacity, or enough time has passed
	if rl.bufSize < rl.n || time.Since(rl.bufHead.t) > rl.duration {
		rl.push(time.Now())
		return false
	}

	return true
}

func (rl *RateLimiter) push(t time.Time) {
	if rl.bufSize == 0 || (rl.bufHead == nil || rl.bufTail == nil) {
		rl.bufHead = &timeNode{t: t}
		rl.bufTail = rl.bufHead
		rl.bufSize = 1
		return
	}

	rl.bufTail.next = &timeNode{t: t}
	rl.bufTail = rl.bufTail.next
	rl.bufSize++

	if rl.bufSize > rl.n {
		rl.bufHead = rl.bufHead.next
		rl.bufSize--
	}
}
