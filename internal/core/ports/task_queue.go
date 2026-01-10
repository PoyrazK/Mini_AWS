package ports

import (
	"context"
)

type TaskQueue interface {
	Enqueue(ctx context.Context, queueName string, payload interface{}) error
	Dequeue(ctx context.Context, queueName string) (string, error)
}
