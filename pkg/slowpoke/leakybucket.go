package slowpoke

import "sync"

type LeakyBucket[T any] struct {
	maxItems           int
	queue              []T
	leakyRateInSeconds int64
	mu                 sync.Mutex
}

func NewLeakyBucket[T any](maxItems int, leakyRateInSeconds int64) *LeakyBucket[T] {
	l := LeakyBucket[T]{
		maxItems:           maxItems,
		leakyRateInSeconds: leakyRateInSeconds,
	}

	return &l
}

func (l *LeakyBucket[T]) AddItem(item T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.queue) == l.maxItems {
		return
	}

	l.queue = append(l.queue, item)
}

func (l *LeakyBucket[T]) Leaky() *T {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.queue) == 0 {
		return nil
	}

	i := l.queue[0]

	l.queue = l.queue[1:]

	return &i
}
