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
	queue                                 []int
	ticker                                *time.Ticker
	stopCh                                chan struct{}
	mu                                    sync.Mutex
}

func NewBucket(maxT, interval int) *bucket {
	b := bucket{
		MaxTokens:                             maxT,
		IntervalUntilNewTokenIsAddedInSeconds: interval,
		queue:                                 make([]int, maxT),
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

	if len(b.queue) >= b.MaxTokens {
		return
	}

	b.queue = append(b.queue, 1)
	fmt.Printf("Token added. Current tokens: %d at %s\n", len(b.queue), time.Now().Format("15:04:05"))
}

func (b *bucket) GetToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.queue) <= 0 {
		return false
	}

	b.queue = b.queue[:len(b.queue)-1]

	return true
}

func (b *bucket) Stop() {
	close(b.stopCh)
}
