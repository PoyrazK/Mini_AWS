package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

const (
	miniredisStartErrMsg = "failed to start miniredis: %v"
	redisTestQueue       = "test_queue"
)

func TestRedisTaskQueueEnqueue(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf(miniredisStartErrMsg, err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	queue := NewRedisTaskQueue(client)

	ctx := context.Background()
	payload := map[string]string{"key": "value"}

	err = queue.Enqueue(ctx, redisTestQueue, payload)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Verify the item is in the queue
	len, err := client.LLen(ctx, redisTestQueue).Result()
	if err != nil {
		t.Fatalf("LLen failed: %v", err)
	}
	if len != 1 {
		t.Fatalf("expected queue length 1, got %d", len)
	}
}

func TestRedisTaskQueueDequeue(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf(miniredisStartErrMsg, err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	queue := NewRedisTaskQueue(client)

	ctx := context.Background()
	payload := map[string]string{"key": "value"}

	// Enqueue first
	err = queue.Enqueue(ctx, redisTestQueue, payload)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Dequeue
	result, err := queue.Dequeue(ctx, redisTestQueue)
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}
	if result == "" {
		t.Fatalf("expected non-empty result")
	}

	// Verify queue is empty
	len, err := client.LLen(ctx, redisTestQueue).Result()
	if err != nil {
		t.Fatalf("LLen failed: %v", err)
	}
	if len != 0 {
		t.Fatalf("expected queue length 0, got %d", len)
	}
}

func TestRedisTaskQueueDequeueEmpty(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf(miniredisStartErrMsg, err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	queue := NewRedisTaskQueue(client)

	ctx := context.Background()

	// Dequeue from empty queue with timeout
	result, err := queue.Dequeue(ctx, redisTestQueue)
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}
	if result != "" {
		t.Fatalf("expected empty result, got %s", result)
	}
}
