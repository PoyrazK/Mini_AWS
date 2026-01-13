package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRedisTaskQueue_Enqueue(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	queue := NewRedisTaskQueue(client)

	ctx := context.Background()
	payload := map[string]string{"key": "value"}

	err = queue.Enqueue(ctx, "test_queue", payload)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Verify the item is in the queue
	len, err := client.LLen(ctx, "test_queue").Result()
	if err != nil {
		t.Fatalf("LLen failed: %v", err)
	}
	if len != 1 {
		t.Fatalf("expected queue length 1, got %d", len)
	}
}

func TestRedisTaskQueue_Dequeue(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	queue := NewRedisTaskQueue(client)

	ctx := context.Background()
	payload := map[string]string{"key": "value"}

	// Enqueue first
	err = queue.Enqueue(ctx, "test_queue", payload)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Dequeue
	result, err := queue.Dequeue(ctx, "test_queue")
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}
	if result == "" {
		t.Fatalf("expected non-empty result")
	}

	// Verify queue is empty
	len, err := client.LLen(ctx, "test_queue").Result()
	if err != nil {
		t.Fatalf("LLen failed: %v", err)
	}
	if len != 0 {
		t.Fatalf("expected queue length 0, got %d", len)
	}
}

func TestRedisTaskQueue_DequeueEmpty(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	queue := NewRedisTaskQueue(client)

	ctx := context.Background()

	// Dequeue from empty queue with timeout
	result, err := queue.Dequeue(ctx, "test_queue")
	if err != nil {
		t.Fatalf("Dequeue failed: %v", err)
	}
	if result != "" {
		t.Fatalf("expected empty result, got %s", result)
	}
}
