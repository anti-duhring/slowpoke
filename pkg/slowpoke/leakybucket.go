package slowpoke

import (
	"sync"
	"time"
)

type LeakyBucket interface {
	CanLeak() bool
}

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

	elapsedSeconds := now.Sub(l.lastLeakTimestamp).Seconds()

	if elapsedSeconds > 0 {
		tokensLeaked := int64(elapsedSeconds) * l.leakyRateInSeconds

		// Remove leaked tokens, but don't go below 0
		l.water = max(0, l.water-tokensLeaked)
		l.lastLeakTimestamp = now
	}

	canAddToken := l.water < l.threshold

	if canAddToken {
		l.water++
	}

	return canAddToken
}
