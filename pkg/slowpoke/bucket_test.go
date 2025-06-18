package slowpoke_test

import (
	"testing"
	"time"

	"github.com/anti-duhring/slowpoke/pkg/slowpoke"
)

func TestRateLimit(t *testing.T) {
	maxTokens := int64(10)
	interval := int64(2)
	bucket := slowpoke.NewBucket(maxTokens, interval)

	// Consumes all avaliable tokens
	for i := range maxTokens {
		if !bucket.HasToken() {
			t.Fatalf("Expected true but got false at iteration %d", i)
		}
	}

	// No more tokens left
	if bucket.HasToken() {
		t.Errorf("Expected false but got true after consuming all tokens")
	}

	// Add 2 more tokens
	time.Sleep(4 * time.Second)

	// Consumes it all
	for i := range 2 {
		if !bucket.HasToken() {
			t.Fatalf("Expected true but got false at iteration %d", i)
		}

	}

	// No more tokens left
	if bucket.HasToken() {
		t.Errorf("Expected false but got true after consuming all tokens")
	}
}
