// Package domain defines core business entities.
package domain

import (
	"io"
	"time"

	"github.com/google/uuid"
)

// Object represents stored object metadata in the storage subsystem.
type Object struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	ARN         string     `json:"arn"`
	Bucket      string     `json:"bucket"`
	Key         string     `json:"key"`
	SizeBytes   int64      `json:"size_bytes"`
	ContentType string     `json:"content_type"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Data        io.Reader  `json:"-"` // Stream for reading/writing
}

type Bucket struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type StorageNode struct {
	ID       string    `json:"id"`
	Address  string    `json:"address"` // host:port
	DataDir  string    `json:"data_dir"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}
