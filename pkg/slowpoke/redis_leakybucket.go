package slowpoke

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLeakyBucket struct {
	threshold          int64
	leakyRateInSeconds int64
	client             *redis.Client
	maxRetries         int
}

type redisLeakyBucketData struct {
	Water             int64     `json:"water"`
	LastLeakTimestamp time.Time `json:"last_leak_timestamp"`
}

func NewRedisLeakyBucket(threshold, leakyRateInSeconds int64, client *redis.Client) *RedisLeakyBucket {
	l := RedisLeakyBucket{
		threshold:          threshold,
		leakyRateInSeconds: leakyRateInSeconds,
		client:             client,
		maxRetries:         1000,
	}

	return &l
}

func (l *RedisLeakyBucket) canAddToken(data *redisLeakyBucketData) bool {
	now := time.Now()
	elapsedSeconds := now.Sub(data.LastLeakTimestamp).Seconds()

	if elapsedSeconds > 0 {
		tokensLeaked := int64(elapsedSeconds) * l.leakyRateInSeconds

		// Remove leaked tokens, but don't go below 0
		data.Water = max(0, data.Water-tokensLeaked)
		data.LastLeakTimestamp = now
	}

	return data.Water < l.threshold
}

func (l *RedisLeakyBucket) Allow(key string, ctx context.Context) (bool, error) {
	allow := true

	h := func(tx *redis.Tx) error {
		var data redisLeakyBucketData

		result := tx.HGetAll(ctx, key)
		if err := result.Err(); err != nil && err != redis.Nil {
			return err
		}

		if result.Val() == nil || len(result.Val()) == 0 {
			data = redisLeakyBucketData{
				Water:             0,
				LastLeakTimestamp: time.Now(),
			}
		} else {
			if err := result.Scan(&data); err != nil {
				return err
			}
		}

		if l.canAddToken(&data) {
			data.Water++
			allow = true
		} else {
			allow = false
		}

		_, err := tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
			return p.HSet(ctx, key, data).Err()
		})
		return err
	}

	for range l.maxRetries {
		err := l.client.Watch(ctx, h, key)
		if err == nil {
			// Tx success
			return allow, nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}

		// Return any other error.
		return false, err
	}

	// max retries
	return false, fmt.Errorf("max retries (%d) exceeded", l.maxRetries)
}
