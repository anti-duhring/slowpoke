package slowpoke

import (
	"fmt"
	"sync"
	"time"
)

type Bucket interface {
	startTick()
	AddToken()
	GetToken() bool
	Stop()
}

type bucket struct {
	MaxTokens                             int
	IntervalUntilNewTokenIsAddedInSeconds int
	tokens                                int
	ticker                                *time.Ticker
	stopCh                                chan struct{}
	mu                                    sync.Mutex
}

func NewBucket(maxT, interval int) *bucket {
	b := bucket{
		MaxTokens:                             maxT,
		IntervalUntilNewTokenIsAddedInSeconds: interval,
		tokens:                                maxT,
		stopCh:                                make(chan struct{}),
	}

	b.startTick()

	return &b
}

func (b *bucket) startTick() {
	ticker := time.NewTicker(time.Duration(b.IntervalUntilNewTokenIsAddedInSeconds) * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				b.AddToken()
			case <-b.stopCh:
				b.ticker.Stop()
				fmt.Println("Ticker stopped")
				return
			}
		}
	}()

	b.ticker = ticker
}

func (b *bucket) AddToken() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.tokens == b.MaxTokens {
		return
	}

	b.tokens += 1
	fmt.Printf("Token added. Current tokens: %d at %s\n", b.tokens, time.Now().Format("15:04:05"))
}

func (b *bucket) GetToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.tokens == 0 {
		return false
	}

	b.tokens -= 1

	return true
}

func (b *bucket) Stop() {
	close(b.stopCh)
}
