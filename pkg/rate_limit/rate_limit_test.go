package rate_limit_test

import (
	"testing"
	"time"

	"github.com/anti-duhring/slowpoke/pkg/rate_limit"
)

func TestRateLimit(t *testing.T) {
	maxTokens := 10
	interval := 2
	bucket := rate_limit.NewBucket(maxTokens, interval)

	if bucket == nil {
		t.Fatal("NewBucket returned nil")
	}
	if bucket.MaxTokens != maxTokens {
		t.Errorf("Expected MaxTokens %d, got %d", maxTokens, bucket.MaxTokens)
	}
	if bucket.IntervalUntilNewTokenIsAddedInSeconds != interval {
		t.Errorf("Expected IntervalUntilNewTokenIsAddedInSeconds %d, got %d", interval, bucket.IntervalUntilNewTokenIsAddedInSeconds)
	}

	// Consumes all avaliable tokens
	for i := range maxTokens {
		hasToken := bucket.GetToken()
		if !hasToken {
			t.Errorf("Expected %v, got %v", true, hasToken)
		}

		t.Log(i, "hasToken", hasToken)
	}

	// No more tokens left
	hasToken := bucket.GetToken()
	if hasToken {
		t.Errorf("Expected %v, got %v", false, hasToken)
	}
	t.Log("hasToken", hasToken)

	// Add 2 more tokens
	time.Sleep(time.Duration(interval*2) * time.Second)

	// Consumes it all
	for i := range 2 {
		hasToken := bucket.GetToken()
		if !hasToken {
			t.Errorf("Expected %v, got %v", true, hasToken)
		}

		t.Log(i, "hasToken", hasToken)
	}

	// No more tokens left
	hasToken = bucket.GetToken()
	if hasToken {
		t.Errorf("Expected %v, got %v", false, hasToken)
	}
	t.Log("hasToken", hasToken)

	bucket.Stop()
}
