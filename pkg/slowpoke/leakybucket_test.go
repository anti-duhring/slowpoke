package slowpoke_test

import (
	"sync"
	"testing"
	"time"

	"github.com/anti-duhring/slowpoke/pkg/slowpoke"
)

func TestCanLeak_BelowThreshold(t *testing.T) {
	threshold := int64(3)
	leakyRate := int64(1000)

	lb := slowpoke.NewLeakyBucket(threshold, leakyRate)

	for i := range int(threshold) - 1 {
		if !lb.CanLeak() {
			t.Fatalf("Expected CanLeak to be true when water is below threshold, but got false at iteration %d", i)
		}
	}
}

func TestCanLeak_AtThreshold(t *testing.T) {
	threshold := int64(3)
	leakyRate := int64(1000)

	lb := slowpoke.NewLeakyBucket(threshold, leakyRate)

	for i := range int(threshold) {
		if !lb.CanLeak() {
			t.Fatalf("Expected CanLeak to be true when filling up to threshold, but got false at iteration %d", i)
		}
	}

	if lb.CanLeak() {
		t.Fatal("Expected CanLeak to be false when water is at threshold, but got true")
	}
}

func TestCanLeak_LeakingOverTime(t *testing.T) {
	threshold := int64(5)
	leakyRate := int64(1)

	lb := slowpoke.NewLeakyBucket(threshold, leakyRate)

	for i := range int(threshold) {
		if !lb.CanLeak() {
			t.Fatalf("Initial fill: Expected CanLeak to be true, but got false at iteration %d", i)
		}
	}

	if lb.CanLeak() {
		t.Fatal("Expected CanLeak to be false when full, but got true before time passed")
	}

	time.Sleep(1 * time.Second)

	if !lb.CanLeak() {
		t.Fatal("Expected CanLeak to be true after 1 second leak, but got false")
	}

	time.Sleep(2 * time.Second)

	if !lb.CanLeak() {
		t.Fatal("Expected CanLeak to be true after 2 second leak, but got false")
	}

	if !lb.CanLeak() {
		t.Fatal("Expected CanLeak to be true after 2 second leak, but got false")
	}

	if lb.CanLeak() {
		t.Fatal("Expected CanLeak to be false when full after 2 second leak, but got true")
	}
}

func TestCanLeak_Concurrency(t *testing.T) {
	threshold := int64(100)
	leakyRate := int64(1000)

	lb := slowpoke.NewLeakyBucket(threshold, leakyRate)

	numGoroutines := 100
	requestsPerGoroutine := 10

	var successfulRequests int64
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestsPerGoroutine {
				// We expect some requests to be rejected once the bucket is full.
				if lb.CanLeak() {
					mu.Lock()
					successfulRequests++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	if successfulRequests > threshold {
		t.Fatalf("Concurrent requests exceeded threshold. Expected at most %d, got %d", threshold, successfulRequests)
	}

	if successfulRequests != threshold {
		t.Fatalf("Expected exactly %d successful requests (bucket filled), got %d", threshold, successfulRequests)
	}
}
