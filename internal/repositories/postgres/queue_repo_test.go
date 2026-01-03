//go:build integration

package postgres

import (
	"testing"
	"time"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresQueueRepository_Integration(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewPostgresQueueRepository(db)
	ctx := setupTestUser(t, db)
	userID := appcontext.UserIDFromContext(ctx)

	cleanDB(t, db)

	t.Run("Create and Get Queue", func(t *testing.T) {
		qID := uuid.New()
		q := &domain.Queue{
			ID:                qID,
			UserID:            userID,
			Name:              "test-queue",
			ARN:               "arn:test",
			VisibilityTimeout: 30,
			RetentionDays:     4,
			MaxMessageSize:    256 * 1024,
			Status:            domain.QueueStatusActive,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err := repo.Create(ctx, q)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, qID, userID)
		require.NoError(t, err)
		assert.Equal(t, q.Name, fetched.Name)

		fetchedByName, err := repo.GetByName(ctx, "test-queue", userID)
		require.NoError(t, err)
		assert.Equal(t, qID, fetchedByName.ID)
	})

	t.Run("Message Operations", func(t *testing.T) {
		qID := uuid.New()
		q := &domain.Queue{
			ID:                qID,
			UserID:            userID,
			Name:              "msg-test-queue",
			ARN:               "arn:test:msg",
			VisibilityTimeout: 1, // 1 second for fast testing
			RetentionDays:     4,
			MaxMessageSize:    256 * 1024,
			Status:            domain.QueueStatusActive,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err := repo.Create(ctx, q)
		require.NoError(t, err)

		// 1. Send Message
		msg, err := repo.SendMessage(ctx, qID, "hello world")
		require.NoError(t, err)
		assert.NotNil(t, msg.ID)

		// 2. Receive Message
		msgs, err := repo.ReceiveMessages(ctx, qID, 1, 1)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		assert.Equal(t, "hello world", msgs[0].Body)
		assert.NotEmpty(t, msgs[0].ReceiptHandle)

		// 3. Receive again (should be empty because of visibility timeout)
		msgsNone, err := repo.ReceiveMessages(ctx, qID, 1, 1)
		require.NoError(t, err)
		assert.Empty(t, msgsNone)

		// 4. Wait for visibility timeout
		time.Sleep(1100 * time.Millisecond)

		// 5. Receive again (should be visible now)
		msgsVisible, err := repo.ReceiveMessages(ctx, qID, 1, 1)
		require.NoError(t, err)
		require.Len(t, msgsVisible, 1)
		assert.Equal(t, 2, msgsVisible[0].ReceivedCount)

		// 6. Delete Message
		err = repo.DeleteMessage(ctx, qID, msgsVisible[0].ReceiptHandle)
		require.NoError(t, err)

		// 7. Receive again (should be empty as deleted)
		msgsDeleted, err := repo.ReceiveMessages(ctx, qID, 1, 1)
		require.NoError(t, err)
		assert.Empty(t, msgsDeleted)
	})

	t.Run("Purge Messages", func(t *testing.T) {
		qID := uuid.New()
		repo.Create(ctx, &domain.Queue{
			ID:                qID,
			UserID:            userID,
			Name:              "purge-test",
			ARN:               "arn:purge",
			VisibilityTimeout: 30,
			RetentionDays:     4,
			MaxMessageSize:    256 * 1024,
			Status:            domain.QueueStatusActive,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		})

		repo.SendMessage(ctx, qID, "m1")
		repo.SendMessage(ctx, qID, "m2")

		affected, err := repo.PurgeMessages(ctx, qID)
		require.NoError(t, err)
		assert.Equal(t, int64(2), affected)

		msgs, _ := repo.ReceiveMessages(ctx, qID, 10, 30)
		assert.Empty(t, msgs)
	})
}
