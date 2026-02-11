package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Get retrieves a value from cache and decodes it into T.
// Returns nil when the key does not exist (cache miss).
// Uses gob encoding to serialize all exported fields regardless of JSON tags.
func Get[T any](ctx context.Context, client *redis.Client, key string) (*T, error) {
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("cache get %q: %w", key, err)
	}

	var result T
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&result); err != nil {
		return nil, fmt.Errorf("cache decode %q: %w", key, err)
	}

	return &result, nil
}

// Set encodes the value and stores it in cache with the given TTL.
// Uses gob encoding to serialize all exported fields regardless of JSON tags.
func Set(ctx context.Context, client *redis.Client, key string, value any, ttl time.Duration) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return fmt.Errorf("cache encode %q: %w", key, err)
	}

	if err := client.Set(ctx, key, buf.Bytes(), ttl).Err(); err != nil {
		return fmt.Errorf("cache set %q: %w", key, err)
	}

	return nil
}

// Delete removes one or more keys from cache.
func Delete(ctx context.Context, client *redis.Client, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if err := client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("cache delete: %w", err)
	}

	return nil
}

// DeleteByPrefix removes all keys matching the given prefix using SCAN.
func DeleteByPrefix(ctx context.Context, client *redis.Client, prefix string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("cache scan %q: %w", prefix, err)
		}

		if len(keys) > 0 {
			if err := client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("cache delete by prefix %q: %w", prefix, err)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
