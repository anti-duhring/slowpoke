package rate_limit_test

import (
	"testing"

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

	bucket.Stop()
}
