package trabago

import (
	"sync"
	"time"
)

// A RateLimiter only allows a code block to be executed N number of times, in a give duration of time.
type RateLimiter struct {
	mu       *sync.Mutex
	n        uint          // number of times
	duration time.Duration // window of time to buffer
	bufHead  *timeNode
	bufTail  *timeNode
	bufSize  uint
}

type timeNode struct {
	t    time.Time
	next *timeNode
}

// NewRateLimiter returns a new rate limiter.
func NewRateLimiter(n uint, duration time.Duration) *RateLimiter {
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
	return rl.isRateLimited(time.Now())
}

func (rl *RateLimiter) ConditionalExecute(f func() interface{}) (interface{}, bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if f == nil || rl.isRateLimited(time.Now()) {
		return nil, false
	}
	return f(), true
}

func (rl *RateLimiter) ConditionalExecuteJob(v interface{}, f func(v interface{}) interface{}) (interface{}, bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if f == nil || rl.isRateLimited(time.Now()) {
		return nil, false
	}
	return f(v), true
}

func (rl *RateLimiter) isRateLimited(t time.Time) bool {
	if rl.bufSize == 0 || rl.bufHead == nil || rl.bufTail == nil {
		rl.bufHead = &timeNode{t: t}
		rl.bufTail = rl.bufHead
		rl.bufSize = 1
		return false
	}

	// allow execution if buffer is below capacity, or enough time has passed
	if t.Sub(rl.bufHead.t) > rl.duration || rl.bufSize < rl.n {
		rl.bufTail.next = &timeNode{t: t}
		rl.bufTail = rl.bufTail.next

		if rl.bufSize == rl.n {
			rl.bufHead = rl.bufHead.next
		} else {
			rl.bufSize++
		}

		return false
	}

	return true
}
