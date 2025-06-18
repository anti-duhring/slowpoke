package slowpoke

import (
	"sync"
	"time"
)

type leakyBucket struct {
	threshold          int64
	water              int64
	leakyRateInSeconds int64
	lastLeakTimestamp  time.Time
	mu                 sync.Mutex
}

func NewLeakyBucket(threshold, leakyRateInSeconds int64) *leakyBucket {
	l := leakyBucket{
		threshold:          threshold,
		leakyRateInSeconds: leakyRateInSeconds,
		lastLeakTimestamp:  time.Now(),
	}

	return &l
}

func (l *leakyBucket) CanLeak() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// now - lastLeakTimestamp = x; X is how many seconds passed since last check
	// x / leakyRateInSeconds = y; y is how many of our time window has passed since then
	// y * threshold = z; z is how much it has been leaked since then.
	// If y is 2 (2 seconds have passed) and threshold is 10, 2 * 10 = 20 requests are considered to have leaked
	// Note: on integer divisions the decimal part is truncated
	leakedAmountSinceLastCheck := ((now.Unix() - l.lastLeakTimestamp.Unix()) / l.leakyRateInSeconds) * l.threshold

	if leakedAmountSinceLastCheck > 0 {
		l.water -= leakedAmountSinceLastCheck
		l.lastLeakTimestamp = now
	}

	if l.water < 0 {
		l.water = 0
	}

	if l.water < l.threshold {
		l.water++
		return true
	}

	return false
}
