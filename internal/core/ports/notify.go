package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
)

type NotifyRepository interface {
	CreateTopic(ctx context.Context, topic *domain.Topic) error
	GetTopicByID(ctx context.Context, id, userID uuid.UUID) (*domain.Topic, error)
	GetTopicByName(ctx context.Context, name string, userID uuid.UUID) (*domain.Topic, error)
	ListTopics(ctx context.Context, userID uuid.UUID) ([]*domain.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error

	CreateSubscription(ctx context.Context, sub *domain.Subscription) error
	GetSubscriptionByID(ctx context.Context, id, userID uuid.UUID) (*domain.Subscription, error)
	ListSubscriptions(ctx context.Context, topicID uuid.UUID) ([]*domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error

	// For message delivery
	SaveMessage(ctx context.Context, msg *domain.NotifyMessage) error
}

type NotifyService interface {
	CreateTopic(ctx context.Context, name string) (*domain.Topic, error)
	ListTopics(ctx context.Context) ([]*domain.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error

	Subscribe(ctx context.Context, topicID uuid.UUID, protocol domain.SubscriptionProtocol, endpoint string) (*domain.Subscription, error)
	ListSubscriptions(ctx context.Context, topicID uuid.UUID) ([]*domain.Subscription, error)
	Unsubscribe(ctx context.Context, id uuid.UUID) error

	Publish(ctx context.Context, topicID uuid.UUID, body string) error
}
