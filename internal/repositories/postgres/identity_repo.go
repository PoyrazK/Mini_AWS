package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/errors"
)

type IdentityRepository struct {
	db DB
}

func NewIdentityRepository(db DB) *IdentityRepository {
	return &IdentityRepository{db: db}
}

func (r *IdentityRepository) CreateAPIKey(ctx context.Context, key *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (id, user_id, key, name, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, key.ID, key.UserID, key.Key, key.Name, key.CreatedAt)
	if err != nil {
		return errors.Wrap(errors.Internal, "failed to create api key", err)
	}
	return nil
}

func (r *IdentityRepository) GetAPIKeyByKey(ctx context.Context, keyStr string) (*domain.APIKey, error) {
	query := `
		SELECT id, user_id, key, name, created_at, last_used
		FROM api_keys
		WHERE key = $1
	`
	var key domain.APIKey
	var lastUsed *time.Time
	err := r.db.QueryRow(ctx, query, keyStr).Scan(
		&key.ID, &key.UserID, &key.Key, &key.Name, &key.CreatedAt, &lastUsed,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(errors.Unauthorized, "invalid api key")
		}
		return nil, errors.Wrap(errors.Internal, "failed to get api key", err)
	}
	if lastUsed != nil {
		key.LastUsed = *lastUsed
	}
	return &key, nil
}
func (r *IdentityRepository) GetAPIKeyByID(ctx context.Context, id uuid.UUID) (*domain.APIKey, error) {
	query := `
		SELECT id, user_id, key, name, created_at, last_used
		FROM api_keys
		WHERE id = $1
	`
	var key domain.APIKey
	var lastUsed *time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(
		&key.ID, &key.UserID, &key.Key, &key.Name, &key.CreatedAt, &lastUsed,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(errors.ObjectNotFound, "api key not found")
		}
		return nil, errors.Wrap(errors.Internal, "failed to get api key", err)
	}
	if lastUsed != nil {
		key.LastUsed = *lastUsed
	}
	return &key, nil
}

func (r *IdentityRepository) ListAPIKeysByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.APIKey, error) {
	query := `
		SELECT id, user_id, key, name, created_at, last_used
		FROM api_keys
		WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to list api keys", err)
	}
	defer rows.Close()

	var keys []*domain.APIKey
	for rows.Next() {
		var key domain.APIKey
		var lastUsed *time.Time
		err := rows.Scan(&key.ID, &key.UserID, &key.Key, &key.Name, &key.CreatedAt, &lastUsed)
		if err != nil {
			return nil, errors.Wrap(errors.Internal, "failed to scan api key", err)
		}
		if lastUsed != nil {
			key.LastUsed = *lastUsed
		}
		keys = append(keys, &key)
	}
	return keys, nil
}

func (r *IdentityRepository) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM api_keys WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(errors.Internal, "failed to delete api key", err)
	}
	return nil
}
