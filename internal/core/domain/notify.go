package domain

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	ARN       string    `json:"arn"` // arn:thecloud:notify:local:{user}:topic/{name}
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubscriptionProtocol string

const (
	ProtocolQueue   SubscriptionProtocol = "queue"
	ProtocolWebhook SubscriptionProtocol = "webhook"
)

type Subscription struct {
	ID        uuid.UUID            `json:"id"`
	UserID    uuid.UUID            `json:"user_id"`
	TopicID   uuid.UUID            `json:"topic_id"`
	Protocol  SubscriptionProtocol `json:"protocol"`
	Endpoint  string               `json:"endpoint"` // arn:thecloud:queue:... or https://...
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type NotifyMessage struct {
	ID        uuid.UUID `json:"id"`
	TopicID   uuid.UUID `json:"topic_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
