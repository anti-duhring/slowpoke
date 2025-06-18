package slowpoke

import (
	"sync"
	"time"
)

type Bucket interface {
	HasToken() bool
}

type bucket struct {
	maxTokens               int64
	addTokenRateInSeconds   int64
	tokens                  int64
	mu                      sync.Mutex
	lastTokenAddedTimeStamp time.Time
}

func NewBucket(maxTokens, addTokenRateInSeconds int64) *bucket {
	b := bucket{
		maxTokens:             maxTokens,
		addTokenRateInSeconds: addTokenRateInSeconds,
		tokens:                maxTokens,
	}

	b.lastTokenAddedTimeStamp = time.Now()

	return &b
}

func (b *bucket) HasToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	tokensAddedSinceLastCheck := (now.Unix() - b.lastTokenAddedTimeStamp.Unix()) / b.addTokenRateInSeconds

	if tokensAddedSinceLastCheck > 0 {
		b.tokens += tokensAddedSinceLastCheck
		b.lastTokenAddedTimeStamp = now
	}

	if b.tokens > b.maxTokens {
		b.tokens = b.maxTokens
	}

	if b.tokens == 0 {
		return false
	}

	b.tokens -= 1

	return true
}
